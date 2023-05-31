package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newWhoamiDataSource() datasource.DataSource {
	return &whoamiDataSource{}
}

type whoamiDataSourceConfig struct {
	ID     types.String `tfsdk:"id"`
	Email  types.String `tfsdk:"email"`
	Issuer types.String `tfsdk:"issuer"`
}

type whoamiDataSource struct {
	client *btpcli.ClientFacade
}

func (gen *whoamiDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_whoami", req.ProviderTypeName)
}

func (gen *whoamiDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	gen.client = req.ProviderData.(*btpcli.ClientFacade)
}

func (gen *whoamiDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Returns information about the logged-in user.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "User ID of the logged-in user.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address of the logged-in user.",
				Computed:            true,
			},
			"issuer": schema.StringAttribute{
				MarkdownDescription: "Name of the token issuer.",
				Computed:            true,
			},
		},
	}
}

func (gen *whoamiDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data whoamiDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := gen.client.GetLoggedInUser()
	if user == nil {
		resp.Diagnostics.AddError("No user found", "")
		return
	}

	data.ID = types.StringValue(user.Username)
	data.Email = types.StringValue(user.Email)
	data.Issuer = types.StringValue(user.Issuer)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
