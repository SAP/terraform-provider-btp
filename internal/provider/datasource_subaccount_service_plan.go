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

func newSubaccountServicePlanDataSource() datasource.DataSource {
	return &subaccountServicePlanDataSource{}
}

type subaccountServicePlanDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	OfferingName types.String `tfsdk:"offering_name"`
	/* OUTPUT */
	Ready             types.Bool   `tfsdk:"ready"`
	Description       types.String `tfsdk:"description"`
	CatalogId         types.String `tfsdk:"catalog_id"`
	CatalogName       types.String `tfsdk:"catalog_name"`
	Free              types.Bool   `tfsdk:"free"`
	Bindable          types.Bool   `tfsdk:"bindable"`
	ServiceOfferingId types.String `tfsdk:"serviceoffering_id"`
	CreatedDate       types.String `tfsdk:"created_date"`
	LastModified      types.String `tfsdk:"last_modified"`
	Shareable         types.Bool   `tfsdk:"shareable"`
}

type subaccountServicePlanDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServicePlanDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_plan", req.ProviderTypeName)
}

func (ds *subaccountServicePlanDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServicePlanDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific service plan such as its name, description, and metadata.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service plan.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("name"), path.MatchRoot("offering_name")),
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service plan.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("offering_name")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"offering_name": schema.StringAttribute{
				MarkdownDescription: "The name of the service offering of the plan.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("name")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service plan is ready.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the service plan.",
				Computed:            true,
			},
			"catalog_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service plan in the service broker catalog.",
				Computed:            true,
			},
			"catalog_name": schema.StringAttribute{
				MarkdownDescription: "The name of the associated service broker catalog.",
				Computed:            true,
			},
			"free": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service plan is free.",
				Computed:            true,
			},
			"bindable": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service plan is bindable.",
				Computed:            true,
			},
			"serviceoffering_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service offering.",
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
			"shareable": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service plan supports instance sharing.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountServicePlanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServicePlanDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cliRes servicemanager.ServicePlanResponseObject
	var err error

	if !data.Id.IsNull() {
		cliRes, _, err = ds.cli.Services.Plan.GetById(ctx, data.SubaccountId.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() && !data.OfferingName.IsNull() {
		cliRes, _, err = ds.cli.Services.Plan.GetByName(ctx, data.SubaccountId.ValueString(), data.Name.ValueString(), data.OfferingName.ValueString())
	} else {
		err = fmt.Errorf("neither offering ID, nor offering Name have been provided")
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Plan (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(cliRes.Id)
	data.Name = types.StringValue(cliRes.Name)
	data.Ready = types.BoolValue(cliRes.Ready)
	data.Description = types.StringValue(cliRes.Description)
	data.CatalogId = types.StringValue(cliRes.CatalogId)
	data.CatalogName = types.StringValue(cliRes.CatalogName)
	data.Free = types.BoolValue(cliRes.Free)
	data.Bindable = types.BoolValue(cliRes.Bindable)
	data.ServiceOfferingId = types.StringValue(cliRes.ServiceOfferingId)
	data.CreatedDate = timeToValue(cliRes.CreatedAt)
	data.LastModified = timeToValue(cliRes.UpdatedAt)

	if cliRes.Metadata != nil {
		data.Shareable = types.BoolValue(cliRes.Metadata.SupportsInstanceSharing)
	} else {
		data.Shareable = types.BoolValue(false)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
