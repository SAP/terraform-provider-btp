package provider

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_authz"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/xsuaa_trust"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ list.ListResourceWithConfigure = &subaccountRoleCollectionAssignmentListResource{}

type subaccountRoleCollectionAssignmentListResource struct {
	client *btpcli.ClientFacade
}

type subaccountRoleCollectionAssignmentListFilter struct {
	SubaccountID       types.String `tfsdk:"subaccount_id"`
	RoleCollectionName types.String `tfsdk:"role_collection_name"`
	Origin             types.String `tfsdk:"origin"`
}

type subaccountRoleCollectionAssignmentListEntry struct {
	roleCollectionName string
	origin             string
	username           string
	groupName          string
	attributeName      string
	attributeValue     string
}

func NewSubaccountRoleCollectionAssignmentListResource() list.ListResource {
	return &subaccountRoleCollectionAssignmentListResource{}
}

func (r *subaccountRoleCollectionAssignmentListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subaccount_role_collection_assignment"
}

func (r *subaccountRoleCollectionAssignmentListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*btpcli.ClientFacade)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *btpcli.ClientFacade, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *subaccountRoleCollectionAssignmentListResource) ListResourceConfigSchema(
	_ context.Context,
	_ list.ListResourceSchemaRequest,
	resp *list.ListResourceSchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This list resource allows you to discover role collection assignments within a specific BTP subaccount. The optional role collection name and origin restrict discovery to matching assignments.",
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"role_collection_name": schema.StringAttribute{
				MarkdownDescription: "The optional name of the role collection whose assignments should be returned. When omitted, assignments from all role collections in the subaccount are returned.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The optional identity provider origin whose assignments should be returned. When omitted or empty, assignments from all origins are returned.",
				Optional:            true,
			},
		},
	}
}

func (r *subaccountRoleCollectionAssignmentListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var filter subaccountRoleCollectionAssignmentListFilter
	if diags := req.Config.Get(ctx, &filter); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	roleCollectionNames, diags := r.roleCollectionNames(ctx, filter)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	originResolver, diags := r.originResolver(ctx, filter)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	entries, diags := r.assignmentEntries(ctx, filter, roleCollectionNames, originResolver)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for _, entry := range entries {
			result := req.NewListResult(ctx)
			setSubaccountRoleCollectionAssignmentListIdentity(ctx, &result, filter.SubaccountID, entry)

			if req.IncludeResource {
				resourceValue := subaccountRoleCollectionAssignmentListResourceValue(filter.SubaccountID, entry)
				result.Diagnostics.Append(result.Resource.Set(ctx, &resourceValue)...)
			}

			if !push(result) {
				return
			}
		}
	}
}

type subaccountRoleCollectionAssignmentOriginResolver struct {
	filterOrigin         string
	expectedSAMLEntityID string
	trustConfigurations  xsuaa_trust.TrustConfigurationResponseCollectionObject
}

func (r *subaccountRoleCollectionAssignmentListResource) originResolver(
	ctx context.Context,
	filter subaccountRoleCollectionAssignmentListFilter,
) (subaccountRoleCollectionAssignmentOriginResolver, diag.Diagnostics) {
	resolver := subaccountRoleCollectionAssignmentOriginResolver{filterOrigin: filter.Origin.ValueString()}

	if resolver.filterOrigin != "" {
		// Match the managed resource's Read behavior: failure to resolve the trust
		// configuration leaves the direct origin and entity-ID matching fallbacks.
		if trustConfiguration, _, err := r.client.Security.Trust.GetBySubaccount(ctx, filter.SubaccountID.ValueString(), resolver.filterOrigin); err == nil {
			resolver.expectedSAMLEntityID = trustConfiguration.IdentityProvider
		}
		return resolver, nil
	}

	trustConfigurations, _, err := r.client.Security.Trust.ListBySubaccount(ctx, filter.SubaccountID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Role Collection Assignment (Subaccount)",
			fmt.Sprintf("Failed to list trust configurations needed to resolve assignment origins: %s", err),
		)
		return resolver, diags
	}

	resolver.trustConfigurations = trustConfigurations
	return resolver, nil
}

func (r *subaccountRoleCollectionAssignmentListResource) roleCollectionNames(ctx context.Context, filter subaccountRoleCollectionAssignmentListFilter) ([]string, diag.Diagnostics) {
	if !filter.RoleCollectionName.IsNull() {
		return []string{filter.RoleCollectionName.ValueString()}, nil
	}

	roleCollections, _, err := r.client.Security.RoleCollection.ListBySubaccount(ctx, filter.SubaccountID.ValueString())
	if err != nil {
		var diags diag.Diagnostics
		diags.AddError(
			"API Error Reading Resource Role Collection Assignment (Subaccount)",
			fmt.Sprintf("Failed to list role collections: %s", err),
		)
		return nil, diags
	}

	roleCollectionNames := make([]string, 0, len(roleCollections))
	for _, roleCollection := range roleCollections {
		roleCollectionNames = append(roleCollectionNames, roleCollection.Name)
	}

	return roleCollectionNames, nil
}

func (r *subaccountRoleCollectionAssignmentListResource) assignmentEntries(
	ctx context.Context,
	filter subaccountRoleCollectionAssignmentListFilter,
	roleCollectionNames []string,
	originResolver subaccountRoleCollectionAssignmentOriginResolver,
) ([]subaccountRoleCollectionAssignmentListEntry, diag.Diagnostics) {
	var (
		entries []subaccountRoleCollectionAssignmentListEntry
		diags   diag.Diagnostics
	)

	for _, roleCollectionName := range roleCollectionNames {
		users, _, err := r.client.Security.RoleCollection.GetUserAssignmentsBySubaccount(ctx, filter.SubaccountID.ValueString(), roleCollectionName)
		if err != nil {
			diags.AddError(
				"API Error Reading Resource Role Collection Assignment (Subaccount)",
				fmt.Sprintf("Failed to list user assignments for role collection %q: %s", roleCollectionName, err),
			)
			return nil, diags
		}
		entries = append(entries, subaccountRoleCollectionAssignmentUserEntries(users, roleCollectionName, originResolver.filterOrigin)...)

		roleCollection, _, err := r.client.Security.RoleCollection.GetBySubaccountWithAttributeMappings(ctx, filter.SubaccountID.ValueString(), roleCollectionName)
		if err != nil {
			diags.AddError(
				"API Error Reading Resource Role Collection Assignment (Subaccount)",
				fmt.Sprintf("Failed to list group and attribute assignments for role collection %q: %s", roleCollectionName, err),
			)
			return nil, diags
		}
		samlEntries, err := subaccountRoleCollectionAssignmentSAMLEntries(roleCollection.SamlAttributeAssignment, roleCollectionName, originResolver)
		if err != nil {
			diags.AddError(
				"Unable to Resolve Role Collection Assignment Origin (Subaccount)",
				err.Error(),
			)
			return nil, diags
		}
		entries = append(entries, samlEntries...)
	}

	return entries, diags
}

func subaccountRoleCollectionAssignmentUserEntries(users []xsuaa_authz.UserReference, roleCollectionName, origin string) []subaccountRoleCollectionAssignmentListEntry {
	entries := make([]subaccountRoleCollectionAssignmentListEntry, 0, len(users))
	for _, user := range users {
		if origin == "" || originMatches(user.Origin, origin) {
			resultOrigin := origin
			if resultOrigin == "" {
				resultOrigin = user.Origin
			}
			entries = append(entries, subaccountRoleCollectionAssignmentListEntry{
				roleCollectionName: roleCollectionName,
				origin:             resultOrigin,
				username:           user.Username,
			})
		}
	}
	return entries
}

func subaccountRoleCollectionAssignmentSAMLEntries(
	assignments []xsuaa_authz.SamlAttrAssignment,
	roleCollectionName string,
	originResolver subaccountRoleCollectionAssignmentOriginResolver,
) ([]subaccountRoleCollectionAssignmentListEntry, error) {
	entries := make([]subaccountRoleCollectionAssignmentListEntry, 0, len(assignments))
	for _, assignment := range assignments {
		origin, matchesOrigin, err := originResolver.resolveSAMLOrigin(assignment)
		if err != nil {
			return nil, fmt.Errorf("role collection %q, attribute %q with value %q: %w", roleCollectionName, assignment.AttributeName, assignment.AttributeValue, err)
		}
		if !matchesOrigin {
			continue
		}

		entry := subaccountRoleCollectionAssignmentListEntry{roleCollectionName: roleCollectionName, origin: origin}
		if assignment.AttributeName == "Groups" {
			entry.groupName = assignment.AttributeValue
		} else {
			entry.attributeName = assignment.AttributeName
			entry.attributeValue = assignment.AttributeValue
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r subaccountRoleCollectionAssignmentOriginResolver) resolveSAMLOrigin(assignment xsuaa_authz.SamlAttrAssignment) (string, bool, error) {
	if r.filterOrigin != "" {
		matchesOrigin := samlOriginMatches(assignment.IdentityProvider, assignment.SamlEntityId, r.filterOrigin) ||
			(r.expectedSAMLEntityID != "" && assignment.SamlEntityId == r.expectedSAMLEntityID)
		return r.filterOrigin, matchesOrigin, nil
	}

	// Prefer the API display name when it is itself a configured origin key.
	for _, trustConfiguration := range r.trustConfigurations {
		if originMatches(assignment.IdentityProvider, trustConfiguration.OriginKey) {
			return trustConfiguration.OriginKey, true, nil
		}
	}

	candidates := make(map[string]struct{})
	for _, trustConfiguration := range r.trustConfigurations {
		if samlOriginMatches(assignment.IdentityProvider, assignment.SamlEntityId, trustConfiguration.OriginKey) ||
			samlEntityMatchesTrust(assignment.SamlEntityId, trustConfiguration) {
			candidates[trustConfiguration.OriginKey] = struct{}{}
		}
	}

	if len(candidates) == 1 {
		for origin := range candidates {
			return origin, true, nil
		}
	}
	if len(candidates) == 0 {
		return "", false, fmt.Errorf("no trust configuration matches identity provider %q and SAML entity ID %q", assignment.IdentityProvider, assignment.SamlEntityId)
	}
	return "", false, fmt.Errorf("multiple trust configurations match identity provider %q and SAML entity ID %q", assignment.IdentityProvider, assignment.SamlEntityId)
}

func samlEntityMatchesTrust(samlEntityID string, trustConfiguration xsuaa_trust.TrustConfigurationResponseObject) bool {
	samlHost := identityProviderHost(samlEntityID)
	if samlHost == "" {
		return false
	}
	return samlHost == identityProviderHost(trustConfiguration.IdentityProvider) ||
		samlHost == identityProviderHost(trustConfiguration.Domain)
}

func identityProviderHost(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	parsedURL, err := url.Parse(value)
	if err == nil && parsedURL.Hostname() != "" {
		return strings.ToLower(strings.TrimSuffix(parsedURL.Hostname(), "."))
	}

	parsedURL, err = url.Parse("https://" + value)
	if err != nil {
		return ""
	}
	return strings.ToLower(strings.TrimSuffix(parsedURL.Hostname(), "."))
}

func setSubaccountRoleCollectionAssignmentListIdentity(
	ctx context.Context,
	result *list.ListResult,
	subaccountID types.String,
	entry subaccountRoleCollectionAssignmentListEntry,
) {
	result.Identity.SetAttribute(ctx, path.Root("subaccount_id"), subaccountID)
	result.Identity.SetAttribute(ctx, path.Root("role_collection_name"), types.StringValue(entry.roleCollectionName))
	result.Identity.SetAttribute(ctx, path.Root("origin"), types.StringValue(entry.origin))

	if entry.username != "" {
		result.Identity.SetAttribute(ctx, path.Root("user_name"), types.StringValue(entry.username))
	} else if entry.groupName != "" {
		result.Identity.SetAttribute(ctx, path.Root("group_name"), types.StringValue(entry.groupName))
	} else {
		result.Identity.SetAttribute(ctx, path.Root("attribute_name"), types.StringValue(entry.attributeName))
		result.Identity.SetAttribute(ctx, path.Root("attribute_value"), types.StringValue(entry.attributeValue))
	}
}

func subaccountRoleCollectionAssignmentListResourceValue(
	subaccountID types.String,
	entry subaccountRoleCollectionAssignmentListEntry,
) subaccountRoleCollectionAssignmentType {
	resourceValue := subaccountRoleCollectionAssignmentType{
		SubaccountId:       subaccountID,
		RoleCollectionName: types.StringValue(entry.roleCollectionName),
		Username:           stringNullIfEmpty(entry.username),
		Groupname:          stringNullIfEmpty(entry.groupName),
		AttributeName:      stringNullIfEmpty(entry.attributeName),
		AttributeValue:     stringNullIfEmpty(entry.attributeValue),
		Origin:             types.StringValue(entry.origin),
	}

	switch {
	case entry.username != "":
		resourceValue.Id = types.StringValue(fmt.Sprintf("%s,%s,%s,%s", subaccountID.ValueString(), entry.roleCollectionName, entry.username, entry.origin))
	case entry.groupName != "":
		resourceValue.Id = types.StringValue(fmt.Sprintf("%s,%s,group:%s,%s", subaccountID.ValueString(), entry.roleCollectionName, entry.groupName, entry.origin))
	default:
		resourceValue.Id = types.StringValue(fmt.Sprintf("%s,%s,attribute:%s/%s,%s", subaccountID.ValueString(), entry.roleCollectionName, entry.attributeName, entry.attributeValue, entry.origin))
	}

	return resourceValue
}
