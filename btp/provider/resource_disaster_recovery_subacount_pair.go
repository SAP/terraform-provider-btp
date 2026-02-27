package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newDisasterRecoverySubaccountPairResource() resource.Resource {
	return &disasterRecoverySubaccountPairResource{}
}

type disasterRecoverySubaccountPairResource struct {
	cli *btpcli.ClientFacade
}

func (rs *disasterRecoverySubaccountPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_disaster_recovery_subaccount_pair", req.ProviderTypeName)
}

func (rs *disasterRecoverySubaccountPairResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *disasterRecoverySubaccountPairResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Create a subaccount pair for the specified subaccounts.

__Tip:__
You must be assigned to Central Disaster Recovery Administrator in both subaccounts.
Each subaccount can only be paired to one subaccount.
You can create instance pairs and subscription pairs in paired subaccounts.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-btp-multi-region-guide/how-to-create-multi-region-setup-on-btp">`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the first subaccount to pair with.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"paired_subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the second subaccount to pair with.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"pair_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount pair.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the subaccount pair was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "The user who created the subaccount pair.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"globalaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the globalaccount.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

type DisasterRecoverySubaccountPairResourceIdentityModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
}

func (rs *disasterRecoverySubaccountPairResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"subaccount_id": identityschema.StringAttribute{
				RequiredForImport: true,
			},
		},
	}
}

func (rs *disasterRecoverySubaccountPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DisasterRecoverySubaccountPairType

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := rs.cli.DisasterRecovery.SubaccountPair.Get(ctx, data.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error reading subaccount pair", fmt.Sprintf("%s", err))
		return
	}

	data, diags = disasterRecoverySubaccountPairValueFrom(ctx, data.SubaccountId, data.PairedSubaccountId, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	var identity DisasterRecoverySubaccountPairResourceIdentityModel

	diags = req.Identity.Get(ctx, &identity)
	if diags.HasError() {
		identity = DisasterRecoverySubaccountPairResourceIdentityModel{
			SubaccountId: data.SubaccountId,
		}

		diags = resp.Identity.Set(ctx, identity)
		resp.Diagnostics.Append(diags...)
	}
}

func (rs *disasterRecoverySubaccountPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DisasterRecoverySubaccountPairType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.DisasterRecovery.SubaccountPair.Create(ctx, &btpcli.SubaccountPairCreateInput{
		SubaccountId:     plan.SubaccountId.ValueString(),
		WithSubaccountId: plan.PairedSubaccountId.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error creating subaccount pair", fmt.Sprintf("%s", err))
		return
	}

	subaccountPairData, _, err := rs.cli.DisasterRecovery.SubaccountPair.Get(ctx, plan.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error reading subaccount pair after creation", fmt.Sprintf("%s", err))
		return
	}

	data, diags := disasterRecoverySubaccountPairValueFrom(ctx, plan.SubaccountId, plan.PairedSubaccountId, subaccountPairData)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	identity := DisasterRecoverySubaccountPairResourceIdentityModel{
		SubaccountId: plan.SubaccountId,
	}

	diags = resp.Identity.Set(ctx, identity)
	resp.Diagnostics.Append(diags...)
}

func (rs *disasterRecoverySubaccountPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DisasterRecoverySubaccountPairType

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Disaster Recovery Subaccount Pair", "This resource is not supposed to be updated")
	if resp.Diagnostics.HasError() {
		return
	}

}

func (rs *disasterRecoverySubaccountPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DisasterRecoverySubaccountPairType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := rs.cli.DisasterRecovery.SubaccountPair.Delete(ctx, state.SubaccountId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error deleting subaccount pair", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *disasterRecoverySubaccountPairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), types.StringValue(req.ID))...)
		return
	}

	var identityData DisasterRecoverySubaccountPairResourceIdentityModel
	resp.Diagnostics.Append(req.Identity.Get(ctx, &identityData)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount_id"), identityData.SubaccountId)...)
}
