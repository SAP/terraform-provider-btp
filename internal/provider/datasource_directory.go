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

func newDirectoryDataSource() datasource.DataSource {
	return &directoryDataSource{}
}

type directoryDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory", req.ProviderTypeName)
}

func (ds *directoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Get the details about a directory.

__Tip:__
You must be assigned to the global account admin role, or the directory admin if the directory is configured to manage its authorizations.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Details of the user that created the directory.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the directory.",
				Computed:            true,
			},
			"features": schema.SetAttribute{
				ElementType: types.StringType,
				MarkdownDescription: "The features that are enabled for the directory. Possible values are: " +
					"\n\t - `DEFAULT` All directories have the following basic feature enabled: " +
					"\n\t\t 1. Group and filter subaccounts for reports and filters " +
					"\n\t\t 2. Monitor usage and costs on a directory level (costs only available for contracts that use the consumption-based commercial model)" +
					"\n\t\t 3. Set custom properties and tags to the directory for identification and reporting purposes." +
					"\n\t - `ENTITLEMENTS` Allows the assignment of a quota for services and applications to the directory from the global account quota for distribution to the subaccounts under this directory." +
					"\n\t - `AUTHORIZATIONS` Allows the assignment of users as administrators or viewers of this directory. You must apply this feature in combination with the `ENTITLEMENTS` feature.",
				Computed: true,
			},
			"labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "Set of words or phrases assigned to the directory.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name of the directory.",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "The GUID of the directory's parent entity. Typically this is the global account.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the directory. Possible values are: " +
					"\n\t - `OK` The CRUD operation or series of operations completed successfully." +
					"\n\t - `STARTED` CRUD operation on an entity has started." +
					"\n\t - `CREATING` Creating entity operation is in progress." +
					"\n\t - `UPDATING` Updating entity operation is in progress." +
					"\n\t - `MOVING` Moving entity operation is in progress." +
					"\n\t - `PROCESSING` A series of operations related to the entity is in progress." +
					"\n\t - `DELETING` Deleting entity operation is in progress." +
					"\n\t - `PENDING REVIEW` The processing operation has been stopped for reviewing and can be restarted by the operator." +
					"\n\t - `CANCELLED` The operation or processing was canceled by the operator." +
					"\n\t - `CREATION_FAILED` The creation operation failed, and the entity was not created or was created but cannot be used." +
					"\n\t - `UPDATE_FAILED` The update operation failed, and the entity was not updated." +
					"\n\t - `PROCESSING_FAILED` The processing operations failed." +
					"\n\t - `DELETION_FAILED` The delete operation failed, and the entity was not deleted." +
					"\n\t - `MOVE_FAILED` Entity could not be moved to a different location." +
					"\n\t - `MIGRATING` Migrating entity from NEO to CF.",
				Computed: true,
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "Applies only to directories that have the user authorization management feature enabled. The subdomain becomes part of the path used to access the authorization tenant of the directory. Unique within the defined region.",
				Computed:            true,
			},
		},
	}
}

func (ds *directoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.Directory.Get(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Directory", fmt.Sprintf("%s", err))
		return
	}

	data, diags = directoryValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
