package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"
)

func newSubaccountServicePlatformDataSource() datasource.DataSource {
	return &subaccountServicePlatformDataSource{}
}

type subaccountServicePlatformDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`

	/* OUTPUT */
	Ready        types.Bool   `tfsdk:"ready"`
	PlatformType types.String `tfsdk:"type"`
	Description  types.String `tfsdk:"description"`
	CreatedDate  types.String `tfsdk:"created_date"`
	LastModified types.String `tfsdk:"last_modified"`
	Labels       types.Map    `tfsdk:"labels"`
}

type subaccountServicePlatformDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServicePlatformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_platform", req.ProviderTypeName)
}

func (ds *subaccountServicePlatformDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServicePlatformDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific platform that is registered for service consumption in a subaccount by platform id or by platform name. Details include the platform's name, type, and labels.

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
				MarkdownDescription: "The ID of the platform.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the platform.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the platform.",
				Computed:            true,
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the platform is ready for consumption.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the platform.",
				Computed:            true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"labels": schema.MapAttribute{
				ElementType: types.SetType{
					ElemType: types.StringType,
				},
				MarkdownDescription: "Set of words or phrases assigned to the platform.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountServicePlatformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServicePlatformDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cliRes servicemanager.PlatformResponseObject
	var err error

	if !data.Id.IsNull() {
		cliRes, _, err = ds.cli.Services.Platform.GetById(ctx, data.SubaccountId.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() {
		cliRes, _, err = ds.cli.Services.Platform.GetByName(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	} else {
		err = fmt.Errorf("neither platform ID, nor platform Name have been provided")
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Platform (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(cliRes.Id)
	data.Ready = types.BoolValue(cliRes.Ready)
	data.PlatformType = types.StringValue(cliRes.Type_)
	data.Name = types.StringValue(cliRes.Name)
	data.Description = types.StringValue(cliRes.Description)
	data.CreatedDate = timeToValue(cliRes.CreatedAt)
	data.LastModified = timeToValue(cliRes.UpdatedAt)

	data.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, cliRes.Labels)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
