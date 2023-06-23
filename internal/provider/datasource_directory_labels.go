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

func newDirectoryLabelsDataSource() datasource.DataSource {
	return &directoryLabelsDataSource{}
}

type directoryLabelsDataSourceConfig struct {
	/* INPUT */
	DirectoryId types.String `tfsdk:"directory_id"`
	Id          types.String `tfsdk:"id"`
	/* OUTPUT */
	Values types.Map `tfsdk:"values"`
}

type directoryLabelsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryLabelsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_labels", req.ProviderTypeName)
}

func (ds *directoryLabelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryLabelsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all the user-defined labels that are currently assigned to a specific directory.

__Tip:__
You must be assigned to the global account admin or viewer role. These roles assignments are not needed for directories of which you are the directory admin.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				DeprecationMessage:  "Use the `directory_id` attribute instead",
				MarkdownDescription: "The ID of the directory.",
				Computed:            true,
			},
			"values": schema.MapAttribute{
				ElementType:         types.SetType{ElemType: types.StringType},
				Computed:            true,
				MarkdownDescription: "Contains the label values",
			},
		},
	}
}

func (ds *directoryLabelsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryLabelsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.Label.ListByDirectory(ctx, data.DirectoryId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Labels (Directory)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.DirectoryId

	data.Values, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, cliRes.Labels)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
