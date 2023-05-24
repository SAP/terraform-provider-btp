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

__Tips__
You must be assigned to the global account admin role, or the directory admin if the directory is configured to manage its authorizations.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/8ed4a705efa0431b910056c0acdbf377.html`,
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
				ElementType:         types.StringType,
				MarkdownDescription: "The features that are enabled for the directory. Valid values: - DEFAULT: (Mandatory) All directories have the following basic feature enabled. (1) Group and filter subaccounts for reports and filters, (2) monitor usage and costs on a directory level (costs only available for contracts that use the consumption-based commercial model), and (3) set custom properties and tags to the directory for identification and reporting purposes. - ENTITLEMENTS: (Optional) Allows the assignment of a quota for services and applications to the directory from the global account quota for distribution to the subaccounts under this directory.  - AUTHORIZATIONS: (Optional) Allows the assignment of users as administrators or viewers of this directory. You must apply this feature in combination with the ENTITLEMENTS feature. <br/><b>Valid values:</b>  [DEFAULT] [DEFAULT,ENTITLEMENTS] [DEFAULT,ENTITLEMENTS,AUTHORIZATIONS]<br/>",
				Computed:            true,
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
				MarkdownDescription: "The current state of the directory. * <b>STARTED:</b> CRUD operation on an entity has started. * <b>CREATING:</b> Creating entity operation is in progress. * <b>UPDATING:</b> Updating entity operation is in progress. * <b>MOVING:</b> Moving entity operation is in progress. * <b>PROCESSING:</b> A series of operations related to the entity is in progress. * <b>DELETING:</b> Deleting entity operation is in progress. * <b>OK:</b> The CRUD operation or series of operations completed successfully. * <b>PENDING REVIEW:</b> The processing operation has been stopped for reviewing and can be restarted by the operator. * <b>CANCELLED:</b> The operation or processing was canceled by the operator. * <b>CREATION_FAILED:</b> The creation operation failed, and the entity was not created or was created but cannot be used. * <b>UPDATE_FAILED:</b> The update operation failed, and the entity was not updated. * <b>PROCESSING_FAILED:</b> The processing operations failed. * <b>DELETION_FAILED:</b> The delete operation failed, and the entity was not deleted. * <b>MOVE_FAILED:</b> Entity could not be moved to a different location. * <b>MIGRATING:</b> Migrating entity from NEO to CF.",
				Computed:            true,
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
