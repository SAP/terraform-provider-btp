package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AddMeAsSubaccountAdminAction struct {
	cli *btpcli.ClientFacade
}

type AddMeAsSubaccountAdminActionModel struct {
	SubaccountId types.String `tfsdk:"subaccount_id"`
}

var _ action.Action = &AddMeAsSubaccountAdminAction{}

func NewAddMeAsSubaccountAdminAction() action.Action {
	return &AddMeAsSubaccountAdminAction{}
}

func (a *AddMeAsSubaccountAdminAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_add_me_as_subaccount_admin", req.ProviderTypeName)
}

func (a *AddMeAsSubaccountAdminAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Assigns the currently authenticated user as an administrator of the specified subaccount.


__Notes:__
- This action can be used to grant yourself administrator permissions to a subaccount via Terraform in analogy to the btp CLI command "btp update accounts/subaccount <subaccount_id> --add-me-as-admin".
- Be aware that the execution of the action does not result in any changes to the Terraform state. It is recommended to use this action only in exceptional cases.
- For a consistent setup we recommend using the resource ["btp_subaccount_role_collection_assignment"](https://registry.terraform.io/providers/SAP/btp/latest/docs/resources/subaccount_role_collection_assignment).

__Tip:__
You must be assigned to the global account admin role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
		},
	}
}

func (a *AddMeAsSubaccountAdminAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
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

func (a *AddMeAsSubaccountAdminAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data AddMeAsSubaccountAdminActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := a.cli.Accounts.Subaccount.AddMeAsAdmin(ctx, data.SubaccountId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("API Error Adding Current User as Subaccount Admin", fmt.Sprintf("%s", err))
		return
	}
}
