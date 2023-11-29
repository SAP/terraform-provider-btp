package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

var directoryObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":           types.StringType,
		"created_by":   types.StringType,
		"created_date": types.StringType,
		"description":  types.StringType,
		"features": types.SetType{
			ElemType: types.StringType,
		},
		"labels": types.MapType{
			ElemType: types.SetType{
				ElemType: types.StringType,
			},
		},
		"last_modified": types.StringType,
		"name":          types.StringType,
		"parent_id":     types.StringType,
		"subdomain":     types.StringType,
		"state":         types.StringType,
	},
}

func newDirectoriesDataSource() datasource.DataSource {
	return &directoriesDataSource{}
}

type directoriesType struct {
	Id     types.String `tfsdk:"id"`
	Values types.List   `tfsdk:"values"`
}

type directoriesDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directories", req.ProviderTypeName)
}

func (ds *directoriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all the directories in a global account, including the directories in directories.

__Tip:__
You must be assigned to the admin or viewer role of the global account, directory.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `btp_globalaccount` datasource instead",
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: dataSourceDirectorySchemaAttributes,
				},
				MarkdownDescription: "The subaccounts contained in the global account.",
				Computed:            true,
			},
		},
	}
}

func (ds *directoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoriesType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	gaRes, _, err := ds.cli.Accounts.GlobalAccount.GetWithHierarchy(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Directories", fmt.Sprintf("%s", err))
		return
	}

	dirs := getAllDirectories(ctx, resp, gaRes.Children)

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())

	data.Values, diags = types.ListValueFrom(ctx, directoryObjType, dirs)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
