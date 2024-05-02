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

func newSubaccountServiceBrokerDataSource() datasource.DataSource {
	return &subaccountServiceBrokerDataSource{}
}

type subaccountServiceBrokerDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	/* OUTPUT */
	Ready        types.Bool   `tfsdk:"ready"`
	Description  types.String `tfsdk:"description"`
	BrokerUrl    types.String `tfsdk:"broker_url"`
	CreatedDate  types.String `tfsdk:"created_date"`
	LastModified types.String `tfsdk:"last_modified"`
	Labels       types.Map    `tfsdk:"labels"`
}

type subaccountServiceBrokerDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServiceBrokerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_broker", req.ProviderTypeName)
}

func (ds *subaccountServiceBrokerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServiceBrokerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific service broker registered in a subaccount, such as its name, description, labels, and URL.

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
				MarkdownDescription: "The ID of the service broker.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service broker.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service broker is ready.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the service broker.",
				Computed:            true,
			},
			"broker_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the service broker.",
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
				MarkdownDescription: "Set of words or phrases assigned to the service broker.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountServiceBrokerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServiceBrokerDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cliRes servicemanager.ServiceBrokerResponseObject
	var err error

	if !data.Id.IsNull() {
		cliRes, _, err = ds.cli.Services.Broker.GetById(ctx, data.SubaccountId.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() {
		cliRes, _, err = ds.cli.Services.Broker.GetByName(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	} else {
		err = fmt.Errorf("neither broker ID, nor broker Name have been provided")
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Broker (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(cliRes.Id)
	data.Name = types.StringValue(cliRes.Name)
	data.Ready = types.BoolValue(cliRes.Ready)
	data.Description = types.StringValue(cliRes.Description)
	data.BrokerUrl = types.StringValue(cliRes.BrokerUrl)
	data.CreatedDate = timeToValue(cliRes.CreatedAt)
	data.LastModified = timeToValue(cliRes.UpdatedAt)

	data.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, cliRes.Labels)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
