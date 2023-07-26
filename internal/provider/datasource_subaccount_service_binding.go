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

func newSubaccountServiceBindingDataSource() datasource.DataSource {
	return &subaccountServiceBindingDataSource{}
}

type subaccountServiceBindingDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountServiceBindingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccount_service_binding", req.ProviderTypeName)
}

func (ds *subaccountServiceBindingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountServiceBindingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Gets details about a specific service binding, such as its access details. They are included in its 'credentials' property, and typically include access URLs and credentials.`,
		Attributes: map[string]schema.Attribute{
			"subaccount_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the subaccount.",
				Required:            true,
				Validators: []validator.String{
					uuidvalidator.ValidUUID(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service binding.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
					uuidvalidator.ValidUUID(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service binding.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ready": schema.BoolAttribute{
				MarkdownDescription: "Shows whether the service binding is ready.",
				Computed:            true,
			},
			"service_instance_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the service instance associated with the binding.",
				Computed:            true,
			},
			"context": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Contextual data for the resource.",
				Computed:            true,
			},
			"bind_resource": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Contains the resources associated with the binding.",
				Computed:            true,
			},
			"credentials": schema.StringAttribute{
				MarkdownDescription: "The credentials to access the binding.",
				Computed:            true,
				Sensitive:           true,
			},
			"parameters": schema.StringAttribute{
				MarkdownDescription: "The parameters of the service binding as a valid JSON object.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the service binding. Possible values are: \n" +
					getFormattedValueAsTableRow("state", "description") +
					getFormattedValueAsTableRow("---", "---") +
					getFormattedValueAsTableRow("`in progress`", "The operation or processing is in progress") +
					getFormattedValueAsTableRow("`failed`", "The operation or processing failed") +
					getFormattedValueAsTableRow("`succeeded`", "The operation or processing succeeded"),
				Computed: true,
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
				MarkdownDescription: "The set of words or phrases assigned to the binding.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountServiceBindingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountServiceBindingType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var cliRes servicemanager.ServiceBindingResponseObject
	var err error

	if !data.Id.IsNull() {
		cliRes, _, err = ds.cli.Services.Binding.GetById(ctx, data.SubaccountId.ValueString(), data.Id.ValueString())
	} else if !data.Name.IsNull() {
		cliRes, _, err = ds.cli.Services.Binding.GetByName(ctx, data.SubaccountId.ValueString(), data.Name.ValueString())
	} else {
		err = fmt.Errorf("neither binding ID, nor binding Name have been provided")
	}

	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Service Binding (Subaccount)", fmt.Sprintf("%s", err))
		return
	}

	data, diags = subaccountServiceBindingValueFrom(ctx, cliRes)
	data.Parameters = types.StringNull() // the API doesn't return parameters for already created instances
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
