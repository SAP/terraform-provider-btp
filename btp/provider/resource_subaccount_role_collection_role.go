package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newSubaccountRoleCollectionRoleResource() resource.Resource {
	return &subaccountRoleCollectionRoleResource{}
}

type subaccountRoleAssignmentType struct {
	SubaccountId      types.String `tfsdk:"subaccount_id"`
	Name              types.String `tfsdk:"name"`
	ID                types.String `tfsdk:"id"`
	RoleName          types.String `tfsdk:"role_name"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppID types.String `tfsdk:"role_template_app_id"`
}

type subaccountRoleCollectionRoleResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountRoleCollectionRoleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection_role", req.ProviderTypeName)
}

func (rs *subaccountRoleCollectionRoleResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountRoleCollectionRoleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a role collection role in a subaccount. 

This resource is used to assign individual roles to a role collection container.

### Prerequisites
* You must be assigned to the **admin role** of the subaccount.
* A role collection container must already exist. This can be managed via the **btp_subaccount_role_collection_base** resource.

### Conflict of Authority Warning
> [!CAUTION]
> Roles within a collection can be managed either using the **btp_subaccount_role_collection** resource (which manages the entire set of roles) or by using individual **btp_subaccount_role_collection_role** resources in combination with **btp_subaccount_role_collection_base** — **but the two methods cannot be used together**.
> 
> If you use this resource to add a role to a collection that is also managed by a monolithic **btp_subaccount_role_collection** resource, the two resources will fight for control over the roles list, leading to "flapping" during terraform apply and unexpected deletions.

### Further documentation
For more details on role collections and roles, see the [official SAP BTP documentation](https://help.sap.com/docs/btp/sap-business-technology-platform/role-collections-and-roles-in-global-accounts-directories-and-subaccounts).`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role collection.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `subaccount_id`,`name`,`role_name`,`role_template_app_id` and `role_template_name` attributes instead",
				MarkdownDescription: "The combined unique ID of the role collection as used for import operations.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role_name": schema.StringAttribute{
				MarkdownDescription: "The name of the referenced role.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
	}
}

type SubaccountRoleCollectionRoleResourceIdentityModel struct {
	SubaccountId      types.String `tfsdk:"subaccount_id"`
	Name              types.String `tfsdk:"name"`
	RoleName          types.String `tfsdk:"role_name"`
	RoleTemplateName  types.String `tfsdk:"role_template_name"`
	RoleTemplateAppID types.String `tfsdk:"role_template_app_id"`
}

func (rs *subaccountRoleCollectionRoleResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"subaccount_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"role_name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"role_template_name": identityschema.StringAttribute{
				RequiredForImport: true,
			},
			"role_template_app_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (rs *subaccountRoleCollectionRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountRoleAssignmentType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.GetBySubaccount(ctx, state.SubaccountId.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Role Collection", fmt.Sprintf("%s", err))
		return
	}

	found := false
	for _, role := range cliRes.RoleReferences {
		if role.Name == state.RoleName.ValueString() && role.RoleTemplateAppId == state.RoleTemplateAppID.ValueString() {
			found = true
			break
		}
	}

	if !found {
		// If the role is missing (someone deleted it manually), remove from Terraform state
		resp.State.RemoveResource(ctx)
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() {
		// Setting ID of state - required by hashicorps terraform plugin testing framework for Import . See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
		state.ID = types.StringValue(fmt.Sprintf("%s,%s,%s,%s,%s",
			state.SubaccountId.ValueString(),
			state.Name.ValueString(),
			state.RoleName.ValueString(),
			state.RoleTemplateAppID.ValueString(),
			state.RoleTemplateName.ValueString(),
		))
	}

	resp.State.Set(ctx, &state)

	identity := SubaccountRoleCollectionRoleResourceIdentityModel{
		SubaccountId:      types.StringValue(state.SubaccountId.ValueString()),
		Name:              types.StringValue(state.Name.ValueString()),
		RoleName:          types.StringValue(state.RoleName.ValueString()),
		RoleTemplateName:  types.StringValue(state.RoleTemplateName.ValueString()),
		RoleTemplateAppID: types.StringValue(state.RoleTemplateAppID.ValueString()),
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleCollectionRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountRoleAssignmentType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := rs.cli.Security.Role.AddBySubaccount(
		ctx,
		plan.SubaccountId.ValueString(),
		plan.Name.ValueString(),
		plan.RoleName.ValueString(),
		plan.RoleTemplateAppID.ValueString(),
		plan.RoleTemplateName.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("API Error Adding Role", fmt.Sprintf("%s", err))
		return
	}

	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s,%s,%s,%s",
		plan.SubaccountId.ValueString(),
		plan.Name.ValueString(),
		plan.RoleName.ValueString(),
		plan.RoleTemplateAppID.ValueString(),
		plan.RoleTemplateName.ValueString(),
	))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Set data returned by API in identity
	identity := SubaccountRoleCollectionRoleResourceIdentityModel{
		SubaccountId:      types.StringValue(plan.SubaccountId.ValueString()),
		Name:              types.StringValue(plan.Name.ValueString()),
		RoleName:          types.StringValue(plan.RoleName.ValueString()),
		RoleTemplateName:  types.StringValue(plan.RoleTemplateName.ValueString()),
		RoleTemplateAppID: types.StringValue(plan.RoleTemplateAppID.ValueString()),
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleCollectionRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op: Changes to role assignments require replacement (Delete/Create)
	// because the BTP API doesn't support "renaming" a role inside a collection.
}

func (rs *subaccountRoleCollectionRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountRoleAssignmentType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := rs.cli.Security.Role.RemoveBySubaccount(
		ctx,
		state.SubaccountId.ValueString(),
		state.Name.ValueString(),
		state.RoleName.ValueString(),
		state.RoleTemplateAppID.ValueString(),
		state.RoleTemplateName.ValueString(),
	)

	if err != nil {
		resp.Diagnostics.AddError("API Error Removing Role Assignment", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountRoleCollectionRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
		idParts := strings.Split(req.ID, ",")

		if len(idParts) != 5 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" || idParts[3] == "" || idParts[4] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Import Identifier",
				fmt.Sprintf("Expected: subaccount_id,collection_name,role_name,app_id,template_name. Got: %q", req.ID),
			)
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_name"), idParts[2])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_template_app_id"), idParts[3])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_template_name"), idParts[4])...)
		return
	}

	var identityData SubaccountRoleCollectionRoleResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), identityData.SubaccountId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), identityData.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_name"), identityData.RoleName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_template_name"), identityData.RoleTemplateName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_template_app_id"), identityData.RoleTemplateAppID)...)
}
