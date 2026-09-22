package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UpdateGlobalAccountAction struct {
	cli *btpcli.ClientFacade
}

type UpdateGlobalAccountActionModel struct {
	DisplayName                   types.String `tfsdk:"display_name"`
	Description                   types.String `tfsdk:"description"`
	EnableSubaccountForceDeletion types.Bool   `tfsdk:"enable_subaccount_force_deletion"`
}

var _ action.Action = &UpdateGlobalAccountAction{}

func NewUpdateGlobalAccountAction() action.Action {
	return &UpdateGlobalAccountAction{}
}

func (a *UpdateGlobalAccountAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_update_global_account", req.ProviderTypeName)
}

func (a *UpdateGlobalAccountAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Updates the settings of the global account.

At least one of ` + "`display_name`" + `, ` + "`description`" + `, or ` + "`enable_subaccount_force_deletion`" + ` must be specified.

__Tip:__
You must be assigned to the global account admin role.

__Further documentation:__
<https://help.sap.com/docs/btp/btp-cli-command-reference/btp-update-accounts-global-account>`,
		Attributes: map[string]schema.Attribute{
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The new display name of the global account.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The new description of the global account.",
				Optional:            true,
			},
			"enable_subaccount_force_deletion": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables the force-deletion of subaccounts in the global account.",
				Optional:            true,
			},
		},
	}
}

func (a *UpdateGlobalAccountAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

func (a *UpdateGlobalAccountAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data UpdateGlobalAccountActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.DisplayName.IsNull() && data.Description.IsNull() && data.EnableSubaccountForceDeletion.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"At least one of 'display_name', 'description', or 'enable_subaccount_force_deletion' must be specified.",
		)
		return
	}

	updateInput := &btpcli.GlobalAccountUpdateInput{}

	var updateFields []string

	if !data.DisplayName.IsNull() {
		updateInput.DisplayName = data.DisplayName.ValueString()
		updateFields = append(updateFields, fmt.Sprintf("display name to %q", data.DisplayName.ValueString()))
	}

	if !data.Description.IsNull() {
		updateInput.Description = data.Description.ValueString()
		updateFields = append(updateFields, fmt.Sprintf("description to %q", data.Description.ValueString()))
	}

	if !data.EnableSubaccountForceDeletion.IsNull() {
		v := data.EnableSubaccountForceDeletion.ValueBool()
		updateInput.EnableSubaccountForceDeletion = &v
		updateFields = append(updateFields, fmt.Sprintf("enable_subaccount_force_deletion to %v", v))
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Updating global account: %s...", strings.Join(updateFields, ", ")),
	})

	_, _, err := a.cli.Accounts.GlobalAccount.Update(ctx, updateInput)

	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Global Account", fmt.Sprintf("%s", err))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Global account updated successfully: %s.", strings.Join(updateFields, ", ")),
	})
}
