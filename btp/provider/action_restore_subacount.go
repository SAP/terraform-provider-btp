package provider

import (
	"context"
	"fmt"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
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
		MarkdownDescription: "Restores a subacount that is in pending deletion.",
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

func (a *RestoreSubaccountAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data RestoreSubaccountActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Execute the restore command

	// Execute the polling for the contractState of the subaccount to change

	/*
		done := make(chan bool)
		// Long running API operation in a goroutine
		go func() {
			httpReq, _ := http.NewRequest(
				http.MethodPut,
				"http://example.com/api/do_thing",
				bytes.NewBuffer([]byte(`{"fake": "data"}`)),
			)

			httpResp, err := a.cli.Do(httpReq)
			if err != nil {
				resp.Diagnostics.AddError(
					"HTTP PUT Error",
					"Error updating data. Please report this issue to the provider developers.",
				)
			}
			done <- true
		}()

		ticker := time.NewTicker(10 * time.Second) // Send message back to practitioner every 10 seconds
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// Once this function is called, the message will be displayed in Terraform
				resp.SendProgress(action.InvokeProgressEvent{
					Message: "Waiting for HTTP request to finish...",
				})
			}
		}
	*/
}
