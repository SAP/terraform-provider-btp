package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountTrustConfigurationResource() resource.Resource {
	return &globalaccountTrustConfigurationResource{}
}

type globalaccountTrustConfigurationResource struct {
	cli *btpcli.ClientFacade
}

func (rs *globalaccountTrustConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount_trust_configuration", req.ProviderTypeName)
}

func (rs *globalaccountTrustConfigurationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *globalaccountTrustConfigurationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Establishes trust from a global account to an Identity Authentication tenant.

__Tip:__
You must be assigned to the admin role of the global account.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/trust-and-federation-with-identity-providers>`,
		Attributes: map[string]schema.Attribute{
			"identity_provider": schema.StringAttribute{
				MarkdownDescription: "The name of the Identity Authentication tenant that you want to connect to the global account.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
			"origin": schema.StringAttribute{
				MarkdownDescription: "The origin of the identity provider.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^.{1,27}-platform$`), "must end with '-platform' and not exceed 36 characters"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
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
			"status": schema.StringAttribute{
				MarkdownDescription: "Determines whether the identity provider is currently 'active' or 'inactive'.",
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

func (rs *globalaccountTrustConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state globalaccountTrustConfigurationType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, rawRes, err := rs.cli.Security.Trust.GetByGlobalAccount(ctx, state.Origin.ValueString())
	if err != nil {
		handleReadErrors(ctx, rawRes, cliRes, resp, err, "Resource Trust Configuration (Global Account)")
		return
	}

	state, diags = globalaccountTrustConfigurationFromValue(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountTrustConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan globalaccountTrustConfigurationType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	if !plan.Origin.IsUnknown() {
		origin := plan.Origin.ValueString()
		cliCreateReq.Origin = &origin
	}

	createRes, _, err := rs.cli.Security.Trust.CreateByGlobalAccount(ctx, cliCreateReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Trust Configuration (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	getRes, _, err := rs.cli.Security.Trust.GetByGlobalAccount(ctx, createRes.OriginKey)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Trust Configuration after Creation (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := globalaccountTrustConfigurationFromValue(ctx, getRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountTrustConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan globalaccountTrustConfigurationType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliUpdateReq := btpcli.TrustConfigurationUpdateInput{
		OriginKey: plan.Origin.ValueString(),
	}

	if !plan.Name.IsUnknown() {
		name := plan.Name.ValueString()
		cliUpdateReq.Name = &name
	}

	if !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		cliUpdateReq.Description = &description
	}

	if !plan.Domain.IsUnknown() {
		domain := plan.Domain.ValueString()
		cliUpdateReq.Domain = &domain
	}

	updateRes, _, err := rs.cli.Security.Trust.UpdateByGlobalAccount(ctx, cliUpdateReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Trust Configuration (Global Account)", fmt.Sprintf("%s", err))
		return
	}

	getRes, _, err := rs.cli.Security.Trust.GetByGlobalAccount(ctx, updateRes.OriginKey)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Trust Configuration after Update (Global Account)", fmt.Sprintf("%s", err))
		return
	}
	state, diags := globalaccountTrustConfigurationFromValue(ctx, getRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *globalaccountTrustConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state globalaccountTrustConfigurationType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.Trust.DeleteByGlobalAccount(ctx, state.Origin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Trust Configuration (Global Account)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *globalaccountTrustConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("origin"), req, resp)
}
