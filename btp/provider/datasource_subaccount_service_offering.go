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

func newSubaccountServiceOfferingDataSource() datasource.DataSource {
	return &subaccountServiceOfferingDataSource{}
}

type subaccountServiceOfferingDataSourceConfig struct {
	/* INPUT */
	SubaccountId types.String `tfsdk:"subaccount_id"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	/* OUTPUT */
	Ready                types.Bool   `tfsdk:"ready"`
	Description          types.String `tfsdk:"description"`
	Bindable             types.Bool   `tfsdk:"bindable"`
	InstancesRetrievable types.Bool   `tfsdk:"instances_retrievable"`
	BindingsRetrievable  types.Bool   `tfsdk:"bindings_retrievable"`
	PlanUpdateable       types.Bool   `tfsdk:"plan_updateable"`
	AllowContextUpdates  types.Bool   `tfsdk:"allow_context_updates"`
	Tags                 types.Set    `tfsdk:"tags"`
	BrokerId             types.String `tfsdk:"broker_id"`
	CatalogId            types.String `tfsdk:"catalog_id"`
	CatalogName          types.String `tfsdk:"catalog_name"`
	CreatedDate          types.String `tfsdk:"created_date"`
	LastModified         types.String `tfsdk:"last_modified"`
}

type subaccountServiceOfferingDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServiceOfferingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_offering", req.ProviderTypeName)
}

func (ds *subaccountServiceOfferingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServiceOfferingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific service offering such as its ID, name, description, metadata, and the associated service brokers.

__Tip:__
You must be assigned to the admin or viewer role of the subaccount.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service offering.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service offering.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service offering is ready to be advertised.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the service offering.",
				Computed:            true,
			},
			"bindable": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service offering is bindable.",
				Computed:            true,
			},
			"instances_retrievable": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service instances associated with the service offering can be retrieved.",
				Computed:            true,
			},
			"bindings_retrievable": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the bindings associated with the service offering can be retrieved.",
				Computed:            true,
			},
			"plan_updateable": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the offered plan can be updated.",
				Computed:            true,
			},
			"allow_context_updates": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the context for the service offering can be updated.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The list of tags for the service offering.",
				Computed:            true,
			},
			"broker_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the broker that provides the service plan.",
				Computed:            true,
			},
			"catalog_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service offering as provided by the catalog.",
				Computed:            true,
			},
			"catalog_name": schema.StringAttribute{
				MarkdownDescription: "The catalog name of the service offering.",
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
		},
	}
}

func (ds *subaccountServiceOfferingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServiceOfferingDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cliRes servicemanager.ServiceOfferingResponseObject
	var err error

	if !data.Id.IsNull() {
		cliRes, _, err = ds.cli.Services.Offering.GetById(ctx, data.SubaccountId.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() {
		cliRes, _, err = ds.cli.Services.Offering.GetByName(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	} else {
		err = fmt.Errorf("neither offering ID, nor offering Name have been provided")
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Offering (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data.Id = types.StringValue(cliRes.Id)
	data.Name = types.StringValue(cliRes.Name)
	data.Ready = types.BoolValue(cliRes.Ready)
	data.Description = types.StringValue(cliRes.Description)
	data.Bindable = types.BoolValue(cliRes.Bindable)
	data.InstancesRetrievable = types.BoolValue(cliRes.InstancesRetrievable)
	data.BindingsRetrievable = types.BoolValue(cliRes.BindingsRetrievable)
	data.PlanUpdateable = types.BoolValue(cliRes.PlanUpdateable)
	data.AllowContextUpdates = types.BoolValue(cliRes.AllowContextUpdates)
	data.BrokerId = types.StringValue(cliRes.BrokerId)
	data.CatalogId = types.StringValue(cliRes.CatalogId)
	data.CatalogName = types.StringValue(cliRes.CatalogName)
	data.CreatedDate = timeToValue(cliRes.CreatedAt)
	data.LastModified = timeToValue(cliRes.UpdatedAt)

	data.Tags, diags = types.SetValueFrom(ctx, types.StringType, cliRes.Tags)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
