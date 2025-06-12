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

func newSubaccountEnvironmentInstanceDataSource() datasource.DataSource {
	return &subaccountEnvironmentInstanceDataSource{}
}

type subaccountEnvironmentInstanceDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountEnvironmentInstanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_environment_instance", req.ProviderTypeName)
}

func (ds *subaccountEnvironmentInstanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountEnvironmentInstanceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets the details of a specific environment instance in a subaccount.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the environment instance.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"broker_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the associated environment broker.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"custom_labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				Computed: true,
			},
			"dashboard_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the service dashboard, which is a web-based management user interface for the service instances.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the environment instance.",
				Computed:            true,
			},
			"environment_type": schema.StringAttribute{
				MarkdownDescription: "The type of the environment instance that is used.",
				Computed:            true,
			},
			"labels": schema.StringAttribute{
				MarkdownDescription: "The set of words or phrases assigned to the environment instance.",
				Computed:            true,
			},
			"landscape_label": schema.StringAttribute{
				MarkdownDescription: "The name of the landscape within the logged-in region on which the environment instance is created.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the environment instance.",
				Computed:            true,
			},
			"operation": schema.StringAttribute{
				MarkdownDescription: "An identifier that represents the last operation. This ID is returned by the environment brokers.",
				Computed:            true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "The configuration parameters for the environment instance.",
				Computed:            true,
			},
			"plan_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service plan for the environment instance in the corresponding service broker's catalog.",
				Computed:            true,
			},
			"plan_name": schema.StringAttribute{
				MarkdownDescription: "The name of the service plan for the environment instance in the corresponding service broker's catalog.",
				Computed:            true,
			},
			"platform_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the platform for the environment instance in the corresponding service broker's catalog.",
				Computed:            true,
			},
			"service_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service for the environment instance in the corresponding service broker's catalog.",
				Computed:            true,
			},
			"service_name": schema.StringAttribute{
				MarkdownDescription: "The name of the service for the environment instance in the corresponding service broker's catalog.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the environment instance. Possible values are: \n" +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
					getFormattedValueAsTableRow("`CREATING`", "Creating entity operation is in progress.") +
					getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation operation failed, and the entity was not created or was created but cannot be used.") +
					getFormattedValueAsTableRow("`UPDATING`", "Updating entity operation is in progress.") +
					getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update operation failed, and the entity was not updated.") +
					getFormattedValueAsTableRow("`DELETING`", "Deleting entity operation is in progress.") +
					getFormattedValueAsTableRow("`DELETION_FAILED`", "The delete operation failed, and the entity was not deleted."),
				Computed: true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant that owns the environment instance.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The last provisioning operation on the environment instance. Possible values are: \n" +
					getFormattedValueAsTableRow("type", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`Provision`", "The environment instance is created.") +
					getFormattedValueAsTableRow("`Update`", "The environment instance is changed.") +
					getFormattedValueAsTableRow("`Deprovision`", "The environment instance is deleted."),
				Computed: true,
			},
		},
	}
}

func (ds *subaccountEnvironmentInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountEnvironmentInstanceDataSourceType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.EnvironmentInstance.Get(ctx, data.SubaccountId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Environment Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data, diags = subaccountEnvironmentInstanceDataSourceValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
