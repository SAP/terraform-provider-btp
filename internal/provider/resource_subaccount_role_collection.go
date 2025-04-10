package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountRoleCollectionResource() resource.Resource {
	return &subaccountRoleCollectionResource{}
}

type subaccountRoleCollectionRoleRefType struct {
	Name              types.String `tfsdk:"name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
}

// TODO This predicate is planned to be replaced by letting the subaccountRoleCollectionRoleRefType implement
// TODO	terraform's attr.Value interface and move this code to its Equal method (see also tfutils.SetDifference).
// TODO This will allow to use types.Set instead of a slice for subaccountRoleCollectionType.Roles below.
func saRoleRefIsEqual(roleA, roleB subaccountRoleCollectionRoleRefType) bool {
	return roleA.Name.Equal(roleB.Name) &&
		roleA.RoleTemplateAppId.Equal(roleB.RoleTemplateAppId) &&
		roleA.RoleTemplateName.Equal(roleB.RoleTemplateName)
}

type subaccountRoleCollectionType struct {
	SubaccountId types.String                          `tfsdk:"subaccount_id"`
	Name         types.String                          `tfsdk:"name"`
	Id           types.String                          `tfsdk:"id"`
	Description  types.String                          `tfsdk:"description"`
	Roles        []subaccountRoleCollectionRoleRefType `tfsdk:"roles"`
}

type subaccountRoleCollectionResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountRoleCollectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection", req.ProviderTypeName)
}

func (rs *subaccountRoleCollectionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountRoleCollectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a role collection in a subaccount.

__Tip__
You must be assigned to the admin role of the subaccount.		

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id` and `name` attributes instead",
				MarkdownDescription: "The combined unique ID of the role collection as used for import operations.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the role collection.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"roles": schema.SetNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced role.",
							Required:            true,
						},
						"role_template_name": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced role template.",
							Required:            true,
						},
						"role_template_app_id": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced template app id.",
							Required:            true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

func (rs *subaccountRoleCollectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountRoleCollectionType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Security.RoleCollection.GetBySubaccount(ctx, state.SubaccountId.ValueString(), state.Name.ValueString())
	if err != nil {
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Role Collection (Subaccount)")
		return
	}

	state.Name = types.StringValue(cliRes.Name)
	state.Description = types.StringValue(cliRes.Description)

	if state.Id.IsNull() || state.Id.IsUnknown() {
		// Setting ID of state - required by hashicorps terraform plugin testing framework for Import . See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
		state.Id = types.StringValue(fmt.Sprintf("%s,%s", state.SubaccountId.ValueString(), cliRes.Name))
	}

	state.Roles = []subaccountRoleCollectionRoleRefType{}
	for _, role := range cliRes.RoleReferences {
		state.Roles = append(state.Roles, subaccountRoleCollectionRoleRefType{
			RoleTemplateName:  types.StringValue(role.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
			Name:              types.StringValue(role.Name),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleCollectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountRoleCollectionType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.CreateBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	for _, role := range plan.Roles {
		_, err := rs.cli.Security.Role.AddBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.Name.ValueString(), role.Name.ValueString(), role.RoleTemplateAppId.ValueString(), role.RoleTemplateName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("API Error Adding Role To Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		}
	}

	plan.Name = types.StringValue(cliRes.Name)
	plan.Description = types.StringValue(cliRes.Description)
	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	plan.Id = types.StringValue(fmt.Sprintf("%s,%s", plan.SubaccountId.ValueString(), cliRes.Name))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleCollectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state subaccountRoleCollectionType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var plan subaccountRoleCollectionType
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.RoleCollection.UpdateBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	toBeRemoved := tfutils.SetDifference(state.Roles, plan.Roles, saRoleRefIsEqual)
	for _, role := range toBeRemoved {
		_, err := rs.cli.Security.Role.RemoveBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.Name.ValueString(), role.Name.ValueString(), role.RoleTemplateAppId.ValueString(), role.RoleTemplateName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("API Error Removing Role From Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		}
	}

	toBeAdded := tfutils.SetDifference(plan.Roles, state.Roles, saRoleRefIsEqual)
	for _, role := range toBeAdded {
		_, err := rs.cli.Security.Role.AddBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.Name.ValueString(), role.Name.ValueString(), role.RoleTemplateAppId.ValueString(), role.RoleTemplateName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("API Error Adding Role From Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		}
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.GetBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state.Description = types.StringValue(cliRes.Description)
	state.Roles = []subaccountRoleCollectionRoleRefType{}
	for _, role := range cliRes.RoleReferences {
		state.Roles = append(state.Roles, subaccountRoleCollectionRoleRefType{
			RoleTemplateName:  types.StringValue(role.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
			Name:              types.StringValue(role.Name),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *subaccountRoleCollectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountRoleCollectionType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.RoleCollection.DeleteBySubaccount(ctx, state.SubaccountId.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountRoleCollectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount_id, name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
