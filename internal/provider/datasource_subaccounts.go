package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

var subaccountObjType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":           types.StringType,
		"beta_enabled": types.BoolType,
		"created_by":   types.StringType,
		"created_date": types.StringType,
		"description":  types.StringType,
		"labels": types.MapType{
			ElemType: types.SetType{
				ElemType: types.StringType,
			},
		},
		"last_modified": types.StringType,
		"name":          types.StringType,
		"parent_id":     types.StringType,
		"parent_features": types.SetType{
			ElemType: types.StringType,
		},
		"region": types.StringType,

		"state":          types.StringType,
		"subdomain":      types.StringType,
		"technical_name": types.StringType,
		"usage":          types.StringType,
	},
}

func newSubaccountsDataSource() datasource.DataSource {
	return &subaccountsDataSource{}
}

type subaccountsType struct {
	LabelsFilter types.String `tfsdk:"labels_filter"`
	Values       types.List   `tfsdk:"values"`
}

type subaccountsDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *subaccountsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_subaccounts", req.ProviderTypeName)
}

func (ds *subaccountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *subaccountsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `List all the subaccounts in a global account, including the subaccounts in directories.

__Tip:__
You must be assigned to the admin or viewer role of the global account, directory.`,
		Attributes: map[string]schema.Attribute{
			"labels_filter": schema.StringAttribute{
				MarkdownDescription: "Filters the response based on the labels query.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique ID of the subaccount.",
							Computed:            true,
						},
						"beta_enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the subaccount can use beta services and applications.",
							Computed:            true,
						},
						"created_by": schema.StringAttribute{
							MarkdownDescription: "Details of the user that created the subaccount.",
							Computed:            true,
						},
						"created_date": schema.StringAttribute{
							MarkdownDescription: "The date and time the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "The description of the subaccount.",
							Computed:            true,
						},
						"labels": schema.MapAttribute{
							ElementType: types.SetType{
								ElemType: types.StringType,
							},
							MarkdownDescription: "Set of words or phrases assigned to the subaccount.",
							Computed:            true,
						},
						"last_modified": schema.StringAttribute{
							MarkdownDescription: "The date and time the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "A descriptive name of the subaccount for customer-facing UIs.",
							Computed:            true,
						},
						"parent_id": schema.StringAttribute{
							MarkdownDescription: "The GUID of the subaccount’s parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the GUID of the global account.",
							Computed:            true,
						},
						"parent_features": schema.SetAttribute{
							ElementType:         types.StringType,
							MarkdownDescription: "The features of parent entity of the subaccount.",
							Computed:            true,
						},
						"region": schema.StringAttribute{
							MarkdownDescription: "The region in which the subaccount was created.",
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "The current state of the subaccount. Possible values are: " + // TODO describe the values listed below
								"\n\t - `OK`" +
								"\n\t - `CANCELED`" +
								"\n\t - `CREATING`" +
								"\n\t - `CREATION_FAILED`" +
								"\n\t - `DELETING`" +
								"\n\t - `DELETION_FAILED`" +
								"\n\t - `MIGRATING`" +
								"\n\t - `MIGRATION_FAILED`" +
								"\n\t - `MIGRATED`" +
								"\n\t - `MOVE_FAILED`" +
								"\n\t - `MOVE_TO_OTHER_GA_FAILED`" +
								"\n\t - `MOVING`" +
								"\n\t - `MOVING_TO_OTHER_GA`" +
								"\n\t - `PENDING_REVIEW`" +
								"\n\t - `PROCESSING`" +
								"\n\t - `PROCESSING_FAILED`" +
								"\n\t - `ROLLBACK_MIGRATION_PROCESSING`" +
								"\n\t - `STARTED`" +
								"\n\t - `SUSPENSION_FAILED`" +
								"\n\t - `UPDATE_ACCOUNT_TYPE_FAILED`" +
								"\n\t - `UPDATE_DIRECTORY_TYPE_FAILED`" +
								"\n\t - `UPDATE_FAILED`" +
								"\n\t - `UPDATING`",
							Computed: true,
						},
						"subdomain": schema.StringAttribute{
							MarkdownDescription: "The subdomain that becomes part of the path used to access the authorization tenant of the subaccount. Must be unique within the defined region. Use only letters (a-z), digits (0-9), and hyphens (not at the start or end). Maximum length is 63 characters. Cannot be changed after the subaccount has been created.",
							Computed:            true,
						},
						"usage": schema.StringAttribute{
							MarkdownDescription: "Whether the subaccount is used for production purposes. This flag can help your cloud operator to take appropriate action when handling incidents that are related to mission-critical accounts in production systems. Do not apply for subaccounts that are used for nonproduction purposes, such as development, testing, and demos. Applying this setting this does not modify the subaccount. Possible values are: " +
								"\n\t- `UNSET` Global account or subaccount admin has not set the production-relevancy flag (default value)." +
								"\n\t- `NOT_USED_FOR_PRODUCTION` Subaccount is not used for production purposes." +
								"\n\t- `USED_FOR_PRODUCTION` Subaccount is used for production purposes.",
							Computed: true,
						},
					},
				},
				MarkdownDescription: "The subaccounts contained in the global account.",
				Computed:            true,
			},
		},
	}
}

func (ds *subaccountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data subaccountsType

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var labelsFilter string
	if !data.LabelsFilter.IsUnknown() {
		labelsFilter = data.LabelsFilter.ValueString()
	}

	cliRes, _, err := ds.cli.Accounts.Subaccount.List(ctx, labelsFilter)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Subaccounts", fmt.Sprintf("%s", err))
		return
	}

	subaccountConfigs := []subaccountType{}

	for _, subaccountRes := range cliRes.Value {
		c := subaccountType{
			ID:           types.StringValue(subaccountRes.Guid),
			BetaEnabled:  types.BoolValue(subaccountRes.BetaEnabled),
			CreatedBy:    types.StringValue(subaccountRes.CreatedBy),
			CreatedDate:  timeToValue(subaccountRes.CreatedDate.Time()),
			Description:  types.StringValue(subaccountRes.Description),
			LastModified: timeToValue(subaccountRes.ModifiedDate.Time()),
			Name:         types.StringValue(subaccountRes.DisplayName),
			ParentID:     types.StringValue(subaccountRes.ParentGUID),
			Region:       types.StringValue(subaccountRes.Region),
			State:        types.StringValue(subaccountRes.State),
			Subdomain:    types.StringValue(subaccountRes.Subdomain),
			Usage:        types.StringValue(subaccountRes.UsedForProduction),
		}

		c.Labels, diags = types.MapValueFrom(ctx, types.SetType{ElemType: types.StringType}, subaccountRes.Labels)
		resp.Diagnostics.Append(diags...)

		c.ParentFeatures, diags = types.SetValueFrom(ctx, types.StringType, subaccountRes.ParentFeatures)
		resp.Diagnostics.Append(diags...)

		subaccountConfigs = append(subaccountConfigs, c)
	}

	data.Values, diags = types.ListValueFrom(ctx, subaccountObjType, subaccountConfigs)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
