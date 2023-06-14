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

var regionType attr.Type = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                        types.StringType,
		"name":                      types.StringType,
		"region":                    types.StringType,
		"domain":                    types.StringType,
		"environment":               types.StringType,
		"iaas_provider":             types.StringType,
		"provisioning_service_url":  types.StringType,
		"saas_registry_service_url": types.StringType,
		"supports_trial":            types.BoolType,
	},
}

func newRegionsDataSource() datasource.DataSource {
	return &regionsDataSource{}
}

type regionDataSourceConfig struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Region                 types.String `tfsdk:"region"`
	Domain                 types.String `tfsdk:"domain"`
	Environment            types.String `tfsdk:"environment"`
	IaasProvider           types.String `tfsdk:"iaas_provider"`
	ProvisioningServiceURL types.String `tfsdk:"provisioning_service_url"`
	SaasRegistryServiceURL types.String `tfsdk:"saas_registry_service_url"`
	SupportsTrial          types.Bool   `tfsdk:"supports_trial"`
}

type regionsDataSourceConfig struct {
	Id types.String `tfsdk:"id"`
	/* OUTPUT */
	Values types.List `tfsdk:"values"`
}

type regionsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *regionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_regions", req.ProviderTypeName)
}

func (ds *regionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *regionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Get all the available regions for a global account.

__Tip:__
You must be assigned to the global account admin or viewer role.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				DeprecationMessage:  "Use the `btp_globalaccount` datasource instead",
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Technical name of the data center. Must be unique within the cloud deployment.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Descriptive name of the data center for customer-facing UIs.",
							Computed:            true,
						},
						"region": schema.StringAttribute{
							MarkdownDescription: "The region in which the data center is located.",
							Computed:            true,
						},
						"domain": schema.StringAttribute{
							MarkdownDescription: "The domain of the data center",
							Computed:            true,
						},
						"environment": schema.StringAttribute{
							MarkdownDescription: "The environment that the data center supports. For example: Kubernetes, Cloud Foundry.",
							Computed:            true,
						},
						"iaas_provider": schema.StringAttribute{
							MarkdownDescription: "The infrastructure provider for the data center. Possible values are: " +
								"\n\t - `AWS` Amazon Web Services." +
								"\n\t - `GCP` Google Cloud Platform." +
								"\n\t - `AZURE` Microsoft Azure." +
								"\n\t - `SAP` SAP BTP (Neo)." +
								"\n\t - `ALI` Alibaba Cloud." +
								"\n\t - `IBM` IBM Cloud.",
							Computed: true,
						},
						"provisioning_service_url": schema.StringAttribute{
							MarkdownDescription: "Provisioning service URL.",
							Computed:            true,
						},
						"saas_registry_service_url": schema.StringAttribute{
							MarkdownDescription: "Saas-Registry service URL.",
							Computed:            true,
						},
						"supports_trial": schema.BoolAttribute{
							MarkdownDescription: "Whether the specified datacenter supports trial accounts.",
							Computed:            true,
						},
					},
				},
				MarkdownDescription: "The regions supported by this global account.",
				Computed:            true,
			},
		},
	}
}

func (ds *regionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data regionsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.AvailableRegion.List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Regions", fmt.Sprintf("%s", err))
		return
	}
	regions := []regionDataSourceConfig{}

	for _, regionConf := range cliRes.Datacenters {
		r := regionDataSourceConfig{
			ID:                     types.StringValue(regionConf.Name),
			Name:                   types.StringValue(regionConf.DisplayName),
			Region:                 types.StringValue(regionConf.Region),
			Domain:                 types.StringValue(regionConf.Domain),
			Environment:            types.StringValue(regionConf.Environment),
			IaasProvider:           types.StringValue(regionConf.IaasProvider),
			ProvisioningServiceURL: types.StringValue(regionConf.ProvisioningServiceUrl),
			SaasRegistryServiceURL: types.StringValue(regionConf.SaasRegistryServiceUrl),
			SupportsTrial:          types.BoolValue(regionConf.SupportsTrial),
		}

		regions = append(regions, r)
	}

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())

	data.Values, diags = types.ListValueFrom(ctx, regionType, regions)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
