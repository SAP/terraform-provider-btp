package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const OriginSapDefault = "sap.default"

func newSubaccountTrustConfigurationResource() resource.Resource {
	return &subaccountTrustConfigurationResource{}
}

type subaccountTrustConfigurationResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountTrustConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_trust_configuration", req.ProviderTypeName)
}

func (rs *subaccountTrustConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountTrustConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Establishes trust from a subaccount to an Identity Authentication tenant.

__Tip:__
You must be assigned to the admin role of the subaccount.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/trust-and-federation-with-identity-providers>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"identity_provider": schema.StringAttribute{
				MarkdownDescription: "The name of the Identity Authentication tenant that you want to connect to the subaccount.",
				Required:            true,
				// No validation for the identity provider name, it is validated by the API
				// Needed for handling of sap.default IdP which has no value for this field
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The tenant's domain which should be used for user logon.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name of the trust configuration.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the trust configuration.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"link_text": schema.StringAttribute{
				MarkdownDescription: "Short string that helps users to identify the link for login.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"available_for_user_logon": schema.BoolAttribute{
				MarkdownDescription: "Determines that end users can choose the trust configuration for login. If not set, the trust configuration can remain active, however only application users that explicitly specify the origin key can use if for login.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"auto_create_shadow_users": schema.BoolAttribute{
				MarkdownDescription: "Determines that any user from the tenant can log in. If not set, only the ones who already have a shadow user can log in.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Determines whether the identity provider is currently 'active' or 'inactive'.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("active"),
				Validators: []validator.String{
					stringvalidator.OneOf("active", "inactive"),
				},
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The origin of the identity provider.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				DeprecationMessage:  "Use the `origin` attribute instead",
				MarkdownDescription: "The origin of the identity provider.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The trust type.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol used to establish trust with the identity provider.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the trust configuration can be modified.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (rs *subaccountTrustConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountTrustConfigurationType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Security.Trust.GetBySubaccount(ctx, state.SubaccountId.ValueString(), state.Origin.ValueString())

	if err != nil {
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Trust Configuration (Subaccount)")
		return
	}

	updatedState, diags := subaccountTrustConfigurationFromValue(ctx, cliRes)
	updatedState.SubaccountId = state.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountTrustConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountTrustConfigurationType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Manual check of IdentityProvider field - not possible via schema validation due to handling of sap.default
	// Create only possible for custom IdP -> value for field must be provided
	if plan.IdentityProvider.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("identity_provider"), "Empty Identity Provider", "To create a trust configuration you must provide a non-empty value for the identity provider")
		return
	}

	cliCreateReq := btpcli.TrustConfigurationCreateInput{
		IdentityProvider: plan.IdentityProvider.ValueString(),
	}

	if !plan.Name.IsUnknown() {
		name := plan.Name.ValueString()
		cliCreateReq.Name = &name
	}

	if !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		cliCreateReq.Description = &description
	}

	if !plan.Domain.IsUnknown() {
		domain := plan.Domain.ValueString()
		cliCreateReq.Domain = &domain
	}

	createRes, _, err := rs.cli.Security.Trust.CreateBySubaccount(ctx, plan.SubaccountId.ValueString(), cliCreateReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Trust Configuration (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	availableForUserLogon := plan.AvailableForUserLogon.ValueBool()
	autoCreateShadowUsers := plan.AutoCreateShadowUsers.ValueBool()
	status := plan.Status.ValueString()
	cliUpdateReq := btpcli.TrustConfigurationUpdateInput{
		OriginKey: createRes.OriginKey,
		// TODO: remove repeating domain and idp, see NGPBUG-364505
		IdentityProvider:      &cliCreateReq.IdentityProvider,
		Domain:                cliCreateReq.Domain,
		AvailableForUserLogon: &availableForUserLogon,
		AutoCreateShadowUsers: &autoCreateShadowUsers,
		Status:                &status,
	}

	if !plan.LinkText.IsUnknown() {
		linkText := plan.LinkText.ValueString()
		cliUpdateReq.LinkText = &linkText
	}

	updateRes, _, err := rs.cli.Security.Trust.UpdateBySubaccount(ctx, plan.SubaccountId.ValueString(), cliUpdateReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Trust Configuration after Creation (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := subaccountTrustConfigurationFromValue(ctx, updateRes)
	state.SubaccountId = plan.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountTrustConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountTrustConfigurationType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var state subaccountTrustConfigurationType
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// sap.default and custom IdP must be handled in different ways
	// Manual check for identity provider needed if the origin is not sap.default
	if state.Origin.ValueString() != OriginSapDefault && plan.IdentityProvider.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(path.Root("identity_provider"), "Empty Identity Provider", "To update a trust configuration you must provide a non-empty value for the identity provider")
		return
	}

	idp := plan.IdentityProvider.ValueString()
	availableForUserLogon := plan.AvailableForUserLogon.ValueBool()
	autoCreateShadowUsers := plan.AutoCreateShadowUsers.ValueBool()
	status := plan.Status.ValueString()
	cliUpdateReq := btpcli.TrustConfigurationUpdateInput{
		OriginKey:             plan.Origin.ValueString(),
		IdentityProvider:      &idp,
		AvailableForUserLogon: &availableForUserLogon,
		AutoCreateShadowUsers: &autoCreateShadowUsers,
		Status:                &status,
	}

	if !plan.Domain.IsUnknown() {
		domain := plan.Domain.ValueString()
		cliUpdateReq.Domain = &domain
	}

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		cliUpdateReq.Name = &name
	}

	if !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		cliUpdateReq.Description = &description
	}

	if !plan.LinkText.IsUnknown() {
		linkText := plan.LinkText.ValueString()
		cliUpdateReq.LinkText = &linkText
	}

	updateRes, _, err := rs.cli.Security.Trust.UpdateBySubaccount(ctx, plan.SubaccountId.ValueString(), cliUpdateReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Trust Configuration (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state, diags = subaccountTrustConfigurationFromValue(ctx, updateRes)
	state.SubaccountId = plan.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountTrustConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountTrustConfigurationType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//If the sap.default IdP is managed via Terraform we must skip the explicit deletion
	//It cannot be deleted and it is sufficient to remove it from the state
	if state.Origin.ValueString() == OriginSapDefault {
		resp.Diagnostics.AddWarning("SAP Default cannot be deleted",
			"It is not possible to delete the trust configuration for origin 'sap.default'. "+
				"Skipping the deletion")
		return
	}

	_, _, err := rs.cli.Security.Trust.DeleteBySubaccount(ctx, state.SubaccountId.ValueString(), state.Origin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Trust Configuration (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountTrustConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount_id,origin. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("origin"), idParts[1])...)
}
