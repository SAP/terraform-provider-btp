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

		"state":     types.StringType,
		"subdomain": types.StringType,
		"usage":     types.StringType,
	},
}

func newSubaccountsDataSource() datasource.DataSource {
	return &subaccountsDataSource{}
}

type subaccountsType struct {
	Id           types.String `tfsdk:"id"`
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
		MarkdownDescription: `Gets all the subaccounts in a global account, including the subaccounts in directories.

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
			"id": schema.StringAttribute{ // required by hashicorps terraform plugin testing framework
				DeprecationMessage:  "Use the `btp_globalaccount` datasource instead",
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},
			"values": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the subaccount.",
							Computed:            true,
						},
						"beta_enabled": schema.BoolAttribute{
							MarkdownDescription: "Shows whether the subaccount can use beta services and applications.",
							Computed:            true,
						},
						"created_by": schema.StringAttribute{
							MarkdownDescription: "The details of the user that created the subaccount.",
							Computed:            true,
						},
						"created_date": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
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
							MarkdownDescription: "The set of words or phrases assigned to the subaccount.",
							Computed:            true,
						},
						"last_modified": schema.StringAttribute{
							MarkdownDescription: "The date and time when the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "A descriptive name of the subaccount for customer-facing UIs.",
							Computed:            true,
						},
						"parent_id": schema.StringAttribute{
							MarkdownDescription: "The ID of the subaccount’s parent entity. If the subaccount is located directly in the global account (not in a directory), then this is the ID of the global account.",
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
							MarkdownDescription: "The current state of the subaccount. Possible values are: \n" +
								getFormattedValueAsTableRow("state", "description") +
								getFormattedValueAsTableRow("---", "---") +
								getFormattedValueAsTableRow("`OK`", "The CRUD operation or series of operations completed successfully.") +
								getFormattedValueAsTableRow("`STARTED`", "CRUD operation on the subaccount has started.") +
								getFormattedValueAsTableRow("`CANCELED`", "The operation or processing was canceled by the operator.") +
								getFormattedValueAsTableRow("`PROCESSING`", "A series of operations related to the subaccount are in progress.") +
								getFormattedValueAsTableRow("`PROCESSING_FAILED`", "The processing operations failed.") +
								getFormattedValueAsTableRow("`CREATING`", "Creating the subaccount is in progress.") +
								getFormattedValueAsTableRow("`CREATION_FAILED`", "The creation operation failed, and the subaccount was not created or was created but cannot be used.") +
								getFormattedValueAsTableRow("`UPDATING`", "Updating the subaccount is in progress.") +
								getFormattedValueAsTableRow("`UPDATE_FAILED`", "The update operation failed, and the subaccount was not updated.") +
								getFormattedValueAsTableRow("`UPDATE_DIRECTORY_TYPE_FAILED`", "The update of the directory type failed.") +
								getFormattedValueAsTableRow("`UPDATE_ACCOUNT_TYPE_FAILED`", "The update of the account type failed.") +
								getFormattedValueAsTableRow("`DELETING`", "Deleting the subaccount is in progress.") +
								getFormattedValueAsTableRow("`DELETION_FAILED`", "The deletion of the subaccount failed, and the subaccount was not deleted.") +
								getFormattedValueAsTableRow("`MOVING`", "Moving the subaccount is in progress.") +
								getFormattedValueAsTableRow("`MOVE_FAILED`", "The moving of the subaccount failed.") +
								getFormattedValueAsTableRow("`MOVING_TO_OTHER_GA`", "Moving the subaccount to another global account is in progress.") +
								getFormattedValueAsTableRow("`MOVE_TO_OTHER_GA_FAILED`", "Moving the subaccount to another global account failed.") +
								getFormattedValueAsTableRow("`PENDING_REVIEW`", "The processing operation has been stopped for reviewing and can be restarted by the operator.") +
								getFormattedValueAsTableRow("`MIGRATING`", "Migrating the subaccount from Neo to Cloud Foundry.") +
								getFormattedValueAsTableRow("`MIGRATED`", "The migration of the subaccount completed.") +
								getFormattedValueAsTableRow("`MIGRATION_FAILED`", "The migration of the subaccount failed and the subaccount was not migrated.") +
								getFormattedValueAsTableRow("`ROLLBACK_MIGRATION_PROCESSING`", "The migration of the subaccount was rolled back and the subaccount is not migrated.") +
								getFormattedValueAsTableRow("`SUSPENSION_FAILED`", "The suspension operations failed."),
							Computed: true,
						},
						"subdomain": schema.StringAttribute{
							MarkdownDescription: "The subdomain that becomes part of the path used to access the authorization tenant of the subaccount. Must be unique within the defined region. Use only letters (a-z), digits (0-9), and hyphens (not at the start or end). Maximum length is 63 characters. Cannot be changed after the subaccount has been created.",
							Computed:            true,
						},
						"usage": schema.StringAttribute{
							MarkdownDescription: "Shows whether the subaccount is used for production purposes. This flag can help your cloud operator to take appropriate action when handling incidents that are related to mission-critical accounts in production systems. Do not apply for subaccounts that are used for nonproduction purposes, such as development, testing, and demos. Applying this setting this does not modify the subaccount. Possible values are: \n" +
								getFormattedValueAsTableRow("value", "description") +
								getFormattedValueAsTableRow("---", "---") +
								getFormattedValueAsTableRow("`UNSET`", "Global account or subaccount admin has not set the production-relevancy flag (default value).") +
								getFormattedValueAsTableRow("`NOT_USED_FOR_PRODUCTION`", "The subaccount is not used for production purposes.") +
								getFormattedValueAsTableRow("`USED_FOR_PRODUCTION`", "The subaccount is used for production purposes."),
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

	data.Id = types.StringValue(ds.cli.GetGlobalAccountSubdomain())

	data.Values, diags = types.ListValueFrom(ctx, subaccountObjType, subaccountConfigs)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
