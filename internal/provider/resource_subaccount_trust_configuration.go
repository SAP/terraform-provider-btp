package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

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
		MarkdownDescription: `Establish trust from a subaccount to an Identity Authentication tenant.

__Further documentation:__
https://help.sap.com/docs/BTP/65de2977205c403bbc107264b8eccf4b/cb1bc8f1bd5c482e891063960d7acd78.html`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"identity_provider": schema.StringAttribute{
				MarkdownDescription: "The name of the Identity Authentication tenant that you want the subaccount to connect.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the identity provider.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The origin of the identity provider.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^.{1,27}-platform$`), "must end with '-platform' and not exceed 36 characters"),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description for the identity provider.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The origin of the identity provider.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The trust type.",
				Computed:            true,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol used to establish trust with the identity provider.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Whether the identity provider is currently active or not.",
				Computed:            true,
			},
			"read_only": schema.BoolAttribute{
				MarkdownDescription: "Whether the trust configuration can be modified.",
				Computed:            true,
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

	cliRes, _, err := rs.cli.Security.Trust.GetBySubaccount(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Trust Configuration (Subaccount)", fmt.Sprintf("%s", err))
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

	cliReq := btpcli.TrustConfigurationInput{
		IdentityProvider: plan.IdentityProvider.ValueString(),
	}

	if !plan.Name.IsUnknown() {
		name := plan.Name.ValueString()
		cliReq.Name = &name
	}

	if !plan.Description.IsUnknown() {
		description := plan.Description.ValueString()
		cliReq.Description = &description
	}

	if !plan.Origin.IsUnknown() {
		origin := plan.Origin.ValueString()
		cliReq.Origin = &origin
	}

	createRes, _, err := rs.cli.Security.Trust.CreateBySubaccount(ctx, plan.SubaccountId.ValueString(), cliReq)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Trust Configuration (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	cliRes, _, err := rs.cli.Security.Trust.GetBySubaccount(ctx, plan.SubaccountId.ValueString(), createRes.OriginKey)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Trust Configuration (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	state, diags := subaccountTrustConfigurationFromValue(ctx, cliRes)
	state.SubaccountId = plan.SubaccountId
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountTrustConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountTrustConfigurationType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Trust Configuration (Subaccount)", "This resource is not supposed to be updated")
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *subaccountTrustConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountTrustConfigurationType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.Security.Trust.DeleteBySubaccount(ctx, state.SubaccountId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Trust Configuration (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountTrustConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
