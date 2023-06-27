package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountDataSource() datasource.DataSource {
	return &subaccountDataSource{}
}

type subaccountDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount", req.ProviderTypeName)
}

func (ds *subaccountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a subaccount.

__Tip:__
You must be assigned to the admin or viewer role of the global account, directory, or subaccount.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"beta_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the subaccount can use beta services and applications.",
				Computed:            true,
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Details of the user that created the subaccount.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the subaccount.",
				Computed:            true,
			},
			"labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "Set of words or phrases assigned to the subaccount.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "A descriptive name of the subaccount for customer-facing UIs.",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount’s parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the ID of the global account.",
				Computed:            true,
			},
			"parent_features": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The features of parent entity of the subaccount.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region in which the subaccount was created.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the subaccount. Possible values are: \n" +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
					getFormattedValueAsTableRow("`STARTED`", "CRUD operation on the subaccount has started.") +
					getFormattedValueAsTableRow("`CANCELED`", "The operation or processing was canceled by the operator.") +
					getFormattedValueAsTableRow("`PROCESSING`", "A series of operations related to the subaccount are in progress.") +
					getFormattedValueAsTableRow("`PROCESSING_FAILED`", "The processing operations failed.") +
					getFormattedValueAsTableRow("`CREATING`", "Creating the subaccount is in progress.") +
					getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation operation failed, and the subaccount was not created or was created but cannot be used.") +
					getFormattedValueAsTableRow("`UPDATING`", "Updating the subaccount is in progress.") +
					getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update operation failed, and the subaccount was not updated.") +
					getFormattedValueAsTableRow("`UPDATE_DIRECTORY_TYPE_FAILED`", "The update of the directory type failed.") +
					getFormattedValueAsTableRow("`UPDATE_ACCOUNT_TYPE_FAILED`", "The update of the account type failed.") +
					getFormattedValueAsTableRow("`DELETING`", "Deleting the subaccount is in progress.") +
					getFormattedValueAsTableRow("`DELETION_FAILED`", "The deletion of the subaccount failed, and the subaccount was not deleted.") +
					getFormattedValueAsTableRow("`MOVING`", "Moving the subaccount is in progress.") +
					getFormattedValueAsTableRow("`MOVE_FAILED`", "The moving of the subaccount failed.") +
					getFormattedValueAsTableRow("`MOVING_TO_OTHER_GA`", "Moving the subaccount to another global account is in progress.") +
					getFormattedValueAsTableRow("`MOVE_TO_OTHER_GA_FAILED`", "Moving the subaccount to another global account failed.") +
					getFormattedValueAsTableRow("`PENDING_REVIEW`", "The processing operation has been stopped for reviewing and can be restarted by the operator.") +
					getFormattedValueAsTableRow("`MIGRATING`", "Migrating the subaccount from Neo to Cloud Foundry.") +
					getFormattedValueAsTableRow("`MIGRATED`", "The migration of the subaccount completed.") +
					getFormattedValueAsTableRow("`MIGRATION_FAILED`", "The migration of the subaccount failed and the subaccount was not migrated.") +
					getFormattedValueAsTableRow("`ROLLBACK_MIGRATION_PROCESSING`", "The migration of the subaccount was rolled back and the subaccount is not migrated.") +
					getFormattedValueAsTableRow("`SUSPENSION_FAILED`", "The suspension operations failed."),
				Computed: true,
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "The subdomain that becomes part of the path used to access the authorization tenant of the subaccount. Must be unique within the defined region. Use only letters (a-z), digits (0-9), and hyphens (not at the start or end). Maximum length is 63 characters. Cannot be changed after the subaccount has been created.",
				Computed:            true,
			},
			"usage": schema.StringAttribute{
				MarkdownDescription: "Whether the subaccount is used for production purposes. This flag can help your cloud operator to take appropriate action when handling incidents that are related to mission-critical accounts in production systems. Do not apply for subaccounts that are used for nonproduction purposes, such as development, testing, and demos. Applying this setting this does not modify the subaccount. Possible values are: \n" +
					getFormattedValueAsTableRow("value", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`UNSET`", "Global account or subaccount admin has not set the production-relevancy flag (default value).") +
					getFormattedValueAsTableRow("`NOT_USED_FOR_PRODUCTION`", "The subaccount is not used for production purposes.") +
					getFormattedValueAsTableRow("`USED_FOR_PRODUCTION`", "The subaccount is used for production purposes."),
				Computed: true,
			},
		},
	}
}

func (ds *subaccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.Subaccount.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Subaccount", fmt.Sprintf("%s", err))
		return
	}

	data, diags = subaccountValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
