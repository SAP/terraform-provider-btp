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

func newSubaccountRoleCollectionBaseResource() resource.Resource {
	return &subaccountRoleCollectionBaseResource{}
}

type subaccountRoleCollectionBaseType struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Name         types.String `tfsdk:"name"`
	ID           types.String `tfsdk:"id"`
	Description  types.String `tfsdk:"description"`
}

type subaccountRoleCollectionBaseResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountRoleCollectionBaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_role_collection_base", req.ProviderTypeName)
}

func (rs *subaccountRoleCollectionBaseResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountRoleCollectionBaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Creates a role collection base in a subaccount.

### Prerequisites
* You must be assigned to the **admin role** of the subaccount.

### Conflict of Authority Warning
> [!CAUTION]
> Roles can be defined either directly using the **btp_subaccount_role_collection** resource (which manages the collection and roles together), or by using this **btp_subaccount_role_collection_base** resource in combination with **btp_subaccount_role_collection_role** — **but the two methods cannot be used together**.
> 
> If both the monolithic resource and the individual base/role resources are used against the same Role Collection, spurious changes and conflicting state updates will occur.

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
				DeprecationMessage:  "Use the `subaccount_id` and `name` attributes instead",
				MarkdownDescription: "The combined unique ID of the role collection as used for import operations.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
		},
	}
}

type SubaccountRoleCollectionBaseResourceIdentityModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Name         types.String `tfsdk:"name"`
}

func (rs *subaccountRoleCollectionBaseResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"subaccount_id": identityschema.StringAttribute{
				RequiredForImport: true, // can be defaulted by the provider configuration
			},
			"name": identityschema.StringAttribute{
				RequiredForImport: true, // must be set during import by the practitioner
			},
		},
	}
}

func (rs *subaccountRoleCollectionBaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountRoleCollectionBaseType

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

	if state.ID.IsNull() || state.ID.IsUnknown() {
		// Setting ID of state - required by hashicorps terraform plugin testing framework for Import . See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
		state.ID = types.StringValue(fmt.Sprintf("%s,%s", state.SubaccountId.ValueString(), cliRes.Name))
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var identity SubaccountRoleCollectionBaseResourceIdentityModel

	diags = req.Identity.Get(ctx, &identity)
	if diags.HasError() {
		// During import the identity is not set yet, so set the data returned by API in identity
		identity = SubaccountRoleCollectionBaseResourceIdentityModel{
			SubaccountId: types.StringValue(state.SubaccountId.ValueString()),
			Name:         types.StringValue(cliRes.Name),
		}

		diags = resp.Identity.Set(ctx, identity)
		resp.Diagnostics.Append(diags...)
	}
}

func (rs *subaccountRoleCollectionBaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountRoleCollectionBaseType
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

	plan.Name = types.StringValue(cliRes.Name)
	plan.Description = types.StringValue(cliRes.Description)
	// Setting ID of state - required by hashicorps terraform plugin testing framework for Create. See issue https://github.com/hashicorp/terraform-plugin-testing/issues/84
	plan.ID = types.StringValue(fmt.Sprintf("%s,%s", plan.SubaccountId.ValueString(), cliRes.Name))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Set data returned by API in identity
	identity := SubaccountRoleCollectionBaseResourceIdentityModel{
		SubaccountId: types.StringValue(plan.SubaccountId.ValueString()),
		Name:         types.StringValue(cliRes.Name),
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleCollectionBaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan subaccountRoleCollectionBaseType
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.RoleCollection.UpdateBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	cliRes, _, err := rs.cli.Security.RoleCollection.GetBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Role Collection (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state.Description = types.StringValue(cliRes.Description)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountRoleCollectionBaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountRoleCollectionBaseType
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

func (rs *subaccountRoleCollectionBaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
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
		return
	}

	var identityData SubaccountRoleCollectionBaseResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), identityData.SubaccountId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), identityData.Name)...)
}
