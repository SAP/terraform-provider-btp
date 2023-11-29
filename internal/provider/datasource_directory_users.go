package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newDirectoryUsersDataSource() datasource.DataSource {
	return &directoryUsersDataSource{}
}

type directoryUsersDataSourceConfig struct {
	/* INPUT */
	DirectoryId types.String `tfsdk:"directory_id"`
	Origin      types.String `tfsdk:"origin"`
	Values      types.Set    `tfsdk:"values"`
	Id          types.String `tfsdk:"id"`
}

type directoryUsersDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryUsersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_users", req.ProviderTypeName)
}

func (ds *directoryUsersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryUsersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all users.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/user-and-member-management>`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `directory_id` attribute instead",
				MarkdownDescription: "The ID of the directory.",
				Computed:            true,
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The identity provider that hosts the user. Only needed for custom identity provider..",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"values": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The list of users assigned to the directory.",
				Computed:            true,
			},
		},
	}
}

func (ds *directoryUsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryUsersDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Origin.IsNull() {
		data.Origin = types.StringValue("ldap")
	}

	cliRes, _, err := ds.cli.Security.User.ListByDirectory(ctx, data.DirectoryId.ValueString(), data.Origin.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Users (Directory)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = data.DirectoryId
	data.Values, diags = types.SetValueFrom(ctx, types.StringType, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
