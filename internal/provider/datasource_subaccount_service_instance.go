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

func newSubaccountServiceInstanceDataSource() datasource.DataSource {
	return &subaccountServiceInstanceDataSource{}
}

type subaccountServiceInstanceDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServiceInstanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_instance", req.ProviderTypeName)
}

func (ds *subaccountServiceInstanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServiceInstanceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific provisioned service instance, such as its name, id,  platform to which it belongs, and the last operation performed.

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
				MarkdownDescription: "The ID of the service instance.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service instance.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "The configuration parameters for the service instance.",
				Computed:            true,
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"serviceplan_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service plan.",
				Computed:            true,
			},
			"platform_id": schema.StringAttribute{
				MarkdownDescription: "The platform ID.",
				Computed:            true,
			},
			"referenced_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the instance to which the service instance refers.",
				Computed:            true,
			},
			"shared": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service instance is shared.",
				Computed:            true,
			},
			"context": schema.StringAttribute{
				MarkdownDescription: "Contextual data for the resource.",
				Computed:            true,
			},
			"usable": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the resource can be used.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the service instance.",
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
				MarkdownDescription: "The set of words or phrases assigned to the service instance.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountServiceInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServiceInstanceDataSourceType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cliRes servicemanager.ServiceInstanceResponseObject
	var err error

	if !data.Id.IsNull() {
		cliRes, _, err = ds.cli.Services.Instance.GetById(ctx, data.SubaccountId.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() {
		cliRes, _, err = ds.cli.Services.Instance.GetByName(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	} else {
		err = fmt.Errorf("neither instance ID, nor instance Name have been provided")
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Instance (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data, diags = subaccountServiceInstanceDataSourceValueFrom(ctx, cliRes)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
