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

func newDirectoryEntitlementsDataSource() datasource.DataSource {
	return &directoryEntitlementsDataSource{}
}

type directoryEntitlementsDataSourceConfig struct {
	/* INPUT */
	DirectoryId types.String `tfsdk:"directory_id"`
	Id          types.String `tfsdk:"id"`
	/* OUTPUT */
	Values types.Map `tfsdk:"values"`
}

type directoryEntitlementsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *directoryEntitlementsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_directory_entitlements", req.ProviderTypeName)
}

func (ds *directoryEntitlementsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *directoryEntitlementsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets all the entitlements and quota assignments for a directory.

To view all the resources that a directory and its subdirectories and subaccounts are entitled to use, the following condition must be met:
* The directory must be a directory that is configured to manage its own entitlements.
* You must be assigned to either the global account admin or global account viewers role.`,
		Attributes: map[string]schema.Attribute{
			"directory_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the directory.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{ // required hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `directory_id` attribute instead",
				MarkdownDescription: "The ID of the directory.",
				Computed:            true,
			},
			"values": schema.MapNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service_name": schema.StringAttribute{
							MarkdownDescription: "The name of the entitled service.",
							Computed:            true,
						},
						"service_display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the entitled service.",
							Computed:            true,
						},
						"plan_name": schema.StringAttribute{
							MarkdownDescription: "The name of the entitled service plan.",
							Computed:            true,
						},
						"plan_display_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the entitled service plan.",
							Computed:            true,
						},
						"plan_description": schema.StringAttribute{
							MarkdownDescription: "The description of the entitled service plan.",
							Computed:            true,
						},
						"quota_assigned": schema.Float64Attribute{
							MarkdownDescription: "The overall quota assigned.",
							Computed:            true,
						},
						"quota_remaining": schema.Float64Attribute{
							MarkdownDescription: "The quota, which is not used.",
							Computed:            true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func (ds *directoryEntitlementsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data directoryEntitlementsDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.Entitlement.ListByDirectory(ctx, data.DirectoryId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Entitlements (Directory)", fmt.Sprintf("%s", err))
		return
	}

	values := map[string]entitledService{}

	for _, service := range cliRes.EntitledServices {
		for _, servicePlan := range service.ServicePlans {
			values[fmt.Sprintf("%s:%s", service.Name, servicePlan.Name)] = entitledService{
				ServiceName:        types.StringValue(service.Name),
				ServiceDisplayName: types.StringValue(service.DisplayName),
				PlanName:           types.StringValue(servicePlan.Name),
				PlanDisplayName:    types.StringValue(servicePlan.DisplayName),
				PlanDescription:    types.StringValue(servicePlan.Description),
				QuotaAssigned:      types.Float64Value(servicePlan.Amount),
				QuotaRemaining:     types.Float64Value(servicePlan.RemainingAmount),
			}
		}
	}
	data.Id = data.DirectoryId
	data.Values, diags = types.MapValueFrom(ctx, types.ObjectType{AttrTypes: entitledServiceType()}, values)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
