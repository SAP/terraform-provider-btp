package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/cis"
	"github.com/SAP/terraform-provider-btp/internal/tfutils"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RestoreSubaccountAction struct {
	cli *btpcli.ClientFacade
}

type RestoreSubaccountActionModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
}

var _ action.Action = &RestoreSubaccountAction{}

func NewRestoreSubaccountAction() action.Action {
	return &RestoreSubaccountAction{}
}

func (a *RestoreSubaccountAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_restore_subaccount", req.ProviderTypeName)
}

func (a *RestoreSubaccountAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Cancels the pending deletion of the specified subaccount and restores it to an active state.

__Tip:__
You must be assigned to the global account or directory admin role.

_Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

func (a *RestoreSubaccountAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	cli, ok := req.ProviderData.(*btpcli.ClientFacade)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Action Configure Type",
			fmt.Sprintf("Expected *btpcli.ClientFacade, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	a.cli = cli
}

func (a RestoreSubaccountAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data RestoreSubaccountActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that subaccount with given ID still exists
	cliRes, _, err := a.cli.Accounts.Subaccount.Get(ctx, data.SubaccountId.ValueString())

	if err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("subaccount_id"), "API Error Reading Subaccount", fmt.Sprintf("%s", err))
		return
	}

	if cliRes.ContractStatus != "PENDING_FORCED_DELETION" {
		resp.Diagnostics.AddAttributeError(path.Root("subaccount_id"), "No pending deletion", fmt.Sprintf("The subacount with ID %s is not restorable as it is not in the pending deletion state", data.SubaccountId.ValueString()))
		return
	}

}

func (a *RestoreSubaccountAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RestoreSubaccountActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := a.cli.Accounts.Subaccount.Restore(ctx, data.SubaccountId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("API Error Restoring Resource Subaccount", fmt.Sprintf("%s", err))
		return
	}

	restoreStateConf := &tfutils.StateChangeConf{
		Pending: []string{cis.StateUpdating, cis.StateStarted},
		Target:  []string{cis.StateOK},
		Refresh: func() (any, string, error) {
			subRes, _, err := a.cli.Accounts.Subaccount.Get(ctx, data.SubaccountId.ValueString())

			if err != nil {
				return subRes, "", err
			}

			return subRes, subRes.State, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	restoreRes, err := restoreStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Restoring Resource Subaccount", fmt.Sprintf("%s", err))
	}

	if restoreRes.(cis.SubaccountResponseObject).ContractStatus == "PENDING_FORCED_DELETION" {
		resp.Diagnostics.AddError("API Error Restoring Resource Subaccount", fmt.Sprintf("%s", err))
	}
}
