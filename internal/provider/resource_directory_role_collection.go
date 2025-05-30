package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
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

func newDirectoryRoleCollectionResource() resource.Resource {
	return &directoryRoleCollectionType{}
}

type directoryRoleCollectionRoleRefType struct {
	Name              types.String `tfsdk:"name"`
	RoleTemplateAppId types.String `tfsdk:"role_template_app_id"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
}

// TODO This predicate is planned to be replaced by letting the directoryRoleCollectionRoleRefType implement
// TODO terraform's attr.Value interface and move this code to its Equal method (see also tfutils.SetDifference).
// TODO This will allow to use types.Set instead of a slice for directoryRoleCollectionTypeConfig.Roles below.
func dirRoleRefIsEqual(roleA, roleB directoryRoleCollectionRoleRefType) bool {
	return roleA.Name.Equal(roleB.Name) &&
		roleA.RoleTemplateAppId.Equal(roleB.RoleTemplateAppId) &&
		roleA.RoleTemplateName.Equal(roleB.RoleTemplateName)
}

type directoryRoleCollectionTypeConfig struct {
	Id          types.String                         `tfsdk:"id"`
	DirectoryId types.String                         `tfsdk:"directory_id"`
	Name        types.String                         `tfsdk:"name"`
	Description types.String                         `tfsdk:"description"`
	Roles       []directoryRoleCollectionRoleRefType `tfsdk:"roles"`
}

type directoryRoleCollectionType struct {
	cli *btpcli.ClientFacade
}

func (rs *directoryRoleCollectionType) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_role_collection", req.ProviderTypeName)
}

func (rs *directoryRoleCollectionType) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *directoryRoleCollectionType) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a role collection in a directory.

__Tip:__
You must be assigned to the admin role of the global account or the directory.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts>`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `directory_id` and `name` attributes instead",
				MarkdownDescription: "The combined unique ID of the role collection as used for import operations.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				}},
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
						"role_template_app_id": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced template app id.",
							Required:            true,
						},
						"role_template_name": schema.StringAttribute{
							MarkdownDescription: "The name of the referenced role template.",
							Required:            true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

type DirectoryRoleCollectionResourceIdentityModel struct {
	DirectoryId types.String `tfsdk:"directory_id"`
	Name        types.String `tfsdk:"name"`
}

func (rs *directoryRoleCollectionType) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"directory_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (rs *directoryRoleCollectionType) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state directoryRoleCollectionTypeConfig

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Security.RoleCollection.GetByDirectory(ctx, state.DirectoryId.ValueString(), state.Name.ValueString())
	if err != nil {
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Role Collection (Directory)")
		return
	}

	state.Name = types.StringValue(cliRes.Name)
	state.Description = types.StringValue(cliRes.Description)

	if state.Id.IsNull() || state.Id.IsUnknown() {
		// Setting ID of state - required by hashicorps terraform plugin testing framework for Import. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
		state.Id = types.StringValue(fmt.Sprintf("%s,%s", state.DirectoryId.ValueString(), cliRes.Name))
	}

	state.Roles = []directoryRoleCollectionRoleRefType{}
	for _, role := range cliRes.RoleReferences {
		state.Roles = append(state.Roles, directoryRoleCollectionRoleRefType{
			RoleTemplateName:  types.StringValue(role.RoleTemplateName),
			RoleTemplateAppId: types.StringValue(role.RoleTemplateAppId),
			Name:              types.StringValue(role.Name),
		})
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var identity DirectoryRoleCollectionResourceIdentityModel

	diags = req.Identity.Get(ctx, &identity)
	if diags.HasError() {
		identity = DirectoryRoleCollectionResourceIdentityModel{
			DirectoryId: types.StringValue(state.DirectoryId.ValueString()),
			Name:        types.StringValue(cliRes.Name),
		}

		diags = resp.Identity.Set(ctx, identity)
		resp.Diagnostics.Append(diags...)
	}
}

func (rs *directoryRoleCollectionType) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan directoryRoleCollectionTypeConfig
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.CreateByDirectory(ctx, plan.DirectoryId.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}

	for _, role := range plan.Roles {
		_, err := rs.cli.Security.Role.AddByDirectory(ctx, plan.DirectoryId.ValueString(), plan.Name.ValueString(), role.Name.ValueString(), role.RoleTemplateAppId.ValueString(), role.RoleTemplateName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("API Error Adding Role To Role Collection (Directory)", fmt.Sprintf("%s", err))
		}
	}

	plan.Name = types.StringValue(cliRes.Name)
	plan.Description = types.StringValue(cliRes.Description)

	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	plan.Id = types.StringValue(fmt.Sprintf("%s,%s", plan.DirectoryId.ValueString(), cliRes.Name))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	identity := DirectoryRoleCollectionResourceIdentityModel{
		DirectoryId: types.StringValue(plan.DirectoryId.ValueString()),
		Name:        types.StringValue(cliRes.Name),
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)

}

func (rs *directoryRoleCollectionType) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state directoryRoleCollectionTypeConfig
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var plan directoryRoleCollectionTypeConfig
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.RoleCollection.UpdateByDirectory(ctx, plan.DirectoryId.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}

	toBeRemoved := tfutils.SetDifference(state.Roles, plan.Roles, dirRoleRefIsEqual)
	for _, role := range toBeRemoved {
		_, err := rs.cli.Security.Role.RemoveByDirectory(ctx, plan.DirectoryId.ValueString(), plan.Name.ValueString(), role.Name.ValueString(), role.RoleTemplateAppId.ValueString(), role.RoleTemplateName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("API Error Removing Role From Role Collection (Directory)", fmt.Sprintf("%s", err))
		}
	}

	toBeAdded := tfutils.SetDifference(plan.Roles, state.Roles, dirRoleRefIsEqual)
	for _, role := range toBeAdded {
		_, err := rs.cli.Security.Role.AddByDirectory(ctx, plan.DirectoryId.ValueString(), plan.Name.ValueString(), role.Name.ValueString(), role.RoleTemplateAppId.ValueString(), role.RoleTemplateName.ValueString())

		if err != nil {
			resp.Diagnostics.AddError("API Error Adding Role From Role Collection (Directory)", fmt.Sprintf("%s", err))
		}
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.GetByDirectory(ctx, plan.DirectoryId.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}

	state.Description = types.StringValue(cliRes.Description)
	state.Roles = []directoryRoleCollectionRoleRefType{}
	for _, role := range cliRes.RoleReferences {
		state.Roles = append(state.Roles, directoryRoleCollectionRoleRefType{
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

func (rs *directoryRoleCollectionType) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state directoryRoleCollectionTypeConfig
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.RoleCollection.DeleteByDirectory(ctx, state.DirectoryId.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Role Collection (Directory)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *directoryRoleCollectionType) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
		idParts := strings.Split(req.ID, ",")

		if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected import identifier with format: directory_id, name. Got: %q", req.ID),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_id"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
		return
	}

	var identityData DirectoryRoleCollectionResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("directory_id"), identityData.DirectoryId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), identityData.Name)...)
}
