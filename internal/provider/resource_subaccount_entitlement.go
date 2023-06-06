package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis_entitlements"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountEntitlementResource() resource.Resource {
	return &subaccountEntitlementResource{}
}

type subaccountEntitlementResource struct {
	cli *btpcli.ClientFacade
}

func (rs *subaccountEntitlementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_entitlement", req.ProviderTypeName)
}

func (rs *subaccountEntitlementResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rs.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (rs *subaccountEntitlementResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns the entitlement plan of a service, multitenant application, or environment, to a subaccount. Note that some environments, such as Cloud Foundry, are available by default to all global accounts and their subaccounts, and therefore are not made available as entitlements.

__Tips__
You must be assigned to the global account admin or viewer role.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/00aa2c23479d42568b18882b1ca90d79.html`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the entitled service plan.",
				Computed:            true,
			},
			"service_name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service.",
				Required:            true,
			},
			"plan_name": schema.StringAttribute{
				MarkdownDescription: "The name of the entitled service plan.",
				Required:            true,
			},
			"plan_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the entitled service plan.",
				Computed:            true,
			},
			"amount": schema.Int64Attribute{
				MarkdownDescription: "Quota assigned to the subaccount.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 2000000000),
				},
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the entitlement. Possible values are: " + // TODO describe states listed below
					"\n\t - `OK`" +
					"\n\t - `STARTED`" +
					"\n\t - `PROCESSING`" +
					"\n\t - `PROCESSING_FAILED`",
				Computed: true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
		},
	}
}

func (rs *subaccountEntitlementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state subaccountEntitlementType

	diags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	entitlement, _, err := rs.cli.Accounts.Entitlement.GetAssignedBySubaccount(ctx, state.SubaccountId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Entitlement (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := subaccountEntitlementValueFrom(ctx, *entitlement)
	updatedState.Amount = state.Amount
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountEntitlementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan subaccountEntitlementType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if plan.Amount.IsNull() {
		_, err = rs.cli.Accounts.Entitlement.EnableInSubaccount(ctx, plan.SubaccountId.ValueString(), plan.ServiceName.ValueString(), plan.PlanName.ValueString())
	} else {
		_, err = rs.cli.Accounts.Entitlement.AssignToSubaccount(ctx, plan.SubaccountId.ValueString(), plan.ServiceName.ValueString(), plan.PlanName.ValueString(), int(plan.Amount.ValueInt64()))
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Entitlement (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	// wait for the entitlement to become effective
	createStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis_entitlements.StateStarted, cis_entitlements.StateProcessing},
		Target:  []string{cis_entitlements.StateOK, cis_entitlements.StateProcessingFailed},
		Refresh: func() (interface{}, string, error) {
			entitlement, _, err := rs.cli.Accounts.Entitlement.GetAssignedBySubaccount(ctx, plan.SubaccountId.ValueString(), plan.ServiceName.ValueString(), plan.PlanName.ValueString())

			if err != nil {
				return nil, "", err
			}

			if entitlement == nil {
				return nil, cis_entitlements.StateProcessing, nil
			}

			return *entitlement, entitlement.Assignment.EntityState, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	entitlement, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Resource Entitlement (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	updatedState, diags := subaccountEntitlementValueFrom(ctx, entitlement.(btpcli.UnfoldedEntitlement))
	updatedState.Amount = plan.Amount
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &updatedState)
	resp.Diagnostics.Append(diags...)
}

func (rs *subaccountEntitlementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan subaccountEntitlementType
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("API Error Updating Resource Entitlement (Subaccount)", "Update is not yet implemented.")

	/* TODO: implementation of UPDATE operation
	cliRes, err := gen.client.Execute(ctx, btpcli.Update, gen.command, plan)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Resource Entitlement (Subaccount)", fmt.Sprintf("%s", err))
		return
	}*/

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *subaccountEntitlementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state subaccountEntitlementType
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var err error
	if state.Amount.IsNull() {
		_, err = rs.cli.Accounts.Entitlement.DisableInSubaccount(ctx, state.SubaccountId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString())
	} else {
		_, err = rs.cli.Accounts.Entitlement.AssignToSubaccount(ctx, state.SubaccountId.ValueString(), state.ServiceName.ValueString(), state.PlanName.ValueString(), 0)
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Resource Entitlement (Subaccount)", fmt.Sprintf("%s", err))
		return
	}
}

func (rs *subaccountEntitlementResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: subaccount,service_name,plan_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("subaccount"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plan_name"), idParts[2])...)
}
