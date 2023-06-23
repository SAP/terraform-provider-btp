package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

func newGlobalaccountDataSource() datasource.DataSource {
	return &globalaccountDataSource{}
}

type globalaccountDataSourceConfig struct {
	/* OUTPUT */
	ID               types.String `tfsdk:"id"`
	CommercialModel  types.String `tfsdk:"commercial_model"`
	ConsumptionBased types.Bool   `tfsdk:"consumption_based"`
	ContractStatus   types.String `tfsdk:"contract_status"`
	CostObjectId     types.String `tfsdk:"costobject_id"`
	CostObjectType   types.String `tfsdk:"costobject_type"`
	CreatedDate      types.String `tfsdk:"created_date"`
	CrmCustomerId    types.String `tfsdk:"crm_customer_id"`
	CrmTenantId      types.String `tfsdk:"crm_tenant_id"`
	Description      types.String `tfsdk:"description"`
	DisplayName      types.String `tfsdk:"name"`
	ExpiryDate       types.String `tfsdk:"expiry_date"`
	GeoAccess        types.String `tfsdk:"geo_access"`
	LicenseType      types.String `tfsdk:"license_type"`
	LastModified     types.String `tfsdk:"last_modified"`
	State            types.String `tfsdk:"state"`
	Origin           types.String `tfsdk:"origin"`
	RenewalDate      types.String `tfsdk:"renewal_date"`
	ServiceId        types.String `tfsdk:"service_id"`
	Subdomain        types.String `tfsdk:"subdomain"`
	Usage            types.String `tfsdk:"usage"`
}

type globalaccountDataSource struct {
	cli *btpcli.ClientFacade
}

func (ds *globalaccountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_globalaccount", req.ProviderTypeName)
}

func (ds *globalaccountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ds.cli = req.ProviderData.(*btpcli.ClientFacade)
}

func (ds *globalaccountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Get details about a global account.

__Tip:__
You must be assigned to the global account admin or viewer role.

__Further documentation:__
<https://help.sap.com/docs/btp/sap-business-technology-platform/account-model>`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the global account.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The display name of the global account.",
				Computed:            true,
			},
			"commercial_model": schema.StringAttribute{
				MarkdownDescription: "The type of the commercial contract that was signed.",
				Computed:            true,
			},
			"consumption_based": schema.BoolAttribute{
				MarkdownDescription: "Whether the customer of the global account pays only for services that they actually use (consumption-based) or pay for subscribed services at a fixed cost irrespective of consumption (subscription-based).",
				Computed:            true,
			},
			"contract_status": schema.StringAttribute{
				MarkdownDescription: "The status of the customer contract and its associated root global account. Possible values are: " +
					"\n\t - `ACTIVE` The customer contract and its associated global account is currently active." +
					"\n\t - `PENDING_TERMINATION` A termination process has been triggered for a customer contract (the customer contract has expired, or a customer has given notification that they wish to terminate their contract), and the global account is currently in the validation period. The customer can still access their global account until the end of the validation period." +
					"\n\t - `SUSPENDED` For enterprise accounts, specifies that the customer's global account is currently in the grace period of the termination process. Access to the global account by the customer is blocked. No data is deleted until the deletion date is reached at the end of the grace period. For trial accounts, specifies that the account is suspended, and the account owner has not yet extended the trial period.",
				Computed: true,
			},
			"costobject_id": schema.StringAttribute{
				MarkdownDescription: "The number or code of the cost center, internal order, or Work Breakdown Structure element that is charged for the creation and usage of the global account. The type of the cost object must be configured in `costobject_type`.",
				Computed:            true,
			},
			"costobject_type": schema.StringAttribute{
				MarkdownDescription: "The type of accounting assignment object that is associated with the global account owner and used to charge for the creation and usage of the global account. The number or code of the specified cost object is defined in `costobject_id`. Possible values are: " +
					"\n\t - `COST_CENTER`" +
					"\n\t - `INTERNAL_ORDER`" +
					"\n\t - `WBS_ELEMENT`",
				Computed: true,
			},
			"crm_customer_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the customer as registered in the CRM system.",
				Computed:            true,
			},
			"crm_tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the customer tenant as registered in the CRM system.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the global account.",
				Computed:            true,
			},
			"geo_access": schema.StringAttribute{
				MarkdownDescription: "The geographic locations from where the global account can be accessed. Possible values are: " +
					"\n\t - `STANDARD` The global account can be accessed from any geographic location." +
					"\n\t - `EU_ACCESS` The global account can be accessed only within locations in the EU.",
				Computed: true,
			},
			"license_type": schema.StringAttribute{
				MarkdownDescription: "The type of license for the global account. The license type affects the scope of functions of the account. Possible values are: " +
					"\n\t - `DEVELOPER` For internal developer global accounts on Staging or Canary landscapes." +
					"\n\t - `CUSTOMER` For customer global accounts." +
					"\n\t - `PARTNER` For partner global accounts." +
					"\n\t - `INTERNAL_DEV` For internal global accounts on the Dev landscape." +
					"\n\t - `INTERNAL_PROD` For internal global accounts on the Live landscape." +
					"\n\t - `TRIAL` For customer trial accounts.",
				Computed: true,
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The origin of the account. Possible values are: " +
					"\n\t - `ORDER` Created by the Order Processing API or Submit Order wizard." +
					"\n\t - `OPERATOR` Created by the Global Account wizard." +
					"\n\t - `REGION_SETUP` Created automatically as part of the region setup.",
				Computed: true,
			},
			"service_id": schema.StringAttribute{
				MarkdownDescription: "For internal accounts, the service for which the global account was created.",
				Computed:            true,
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "The subdomain is part of the path used to access the authorization tenant of the global account.",
				Computed:            true,
			},
			"usage": schema.StringAttribute{
				MarkdownDescription: "For internal accounts, the intended purpose of the global account. Possible values are: " +
					"\n\t - `Development` For development of a service." +
					"\n\t - `Testing` For testing development." +
					"\n\t - `Demo` For creating demos." +
					"\n\t - `Production` For delivering a service in a production landscape.",
				Computed: true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the global account. Possible values are: " +
					"\n\t - `STARTED` CRUD operation on an entity has started." +
					"\n\t - `CREATING` Creating entity operation is in progress." +
					"\n\t - `UPDATING` Updating entity operation is in progress." +
					"\n\t - `MOVING` Moving entity operation is in progress." +
					"\n\t - `PROCESSING` A series of operations related to the entity is in progress." +
					"\n\t - `DELETING` Deleting entity operation is in progress." +
					"\n\t - `OK` The CRUD operation or series of operations completed successfully." +
					"\n\t - `PENDING REVIEW` The processing operation has been stopped for reviewing and can be restarted by the operator." +
					"\n\t - `CANCELLED` The operation or processing was canceled by the operator." +
					"\n\t - `CREATION_FAILED` The creation operation failed, and the entity was not created or was created but cannot be used." +
					"\n\t - `UPDATE_FAILED` The update operation failed, and the entity was not updated." +
					"\n\t - `PROCESSING_FAILED` The processing operations failed." +
					"\n\t - `DELETION_FAILED` The delete operation failed, and the entity was not deleted." +
					"\n\t - `MOVE_FAILED` Entity could not be moved to a different location." +
					"\n\t - `MIGRATING` Migrating entity from Neo to Cloud Foundry.",
				Computed: true,
			},
			"created_date": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was created in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"last_modified": schema.StringAttribute{
				MarkdownDescription: "The date and time the resource was last modified in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) format.",
				Computed:            true,
			},
			"expiry_date": schema.StringAttribute{
				MarkdownDescription: "The planned date that the global account expires. This is the same date as the Contract End Date, unless a manual adjustment has been made to the actual expiration date of the global account. Typically, this property is automatically populated only when a formal termination order is received from the CRM system. From a customer perspective, this date marks the start of the grace period, which is typically 30 days before the actual deletion of the account.",
				Computed:            true,
			},
			"renewal_date": schema.StringAttribute{
				MarkdownDescription: "The date that an expired contract was renewed.",
				Computed:            true,
			},
		},
	}
}

func (ds *globalaccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data globalaccountDataSourceConfig

	diags := req.Config.Get(ctx, &data)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cliRes, _, err := ds.cli.Accounts.GlobalAccount.Get(ctx)
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Resource Global Account", fmt.Sprintf("%s", err))
		return
	}

	data.ID = types.StringValue(cliRes.Guid)
	data.CommercialModel = types.StringValue(cliRes.CommercialModel)
	data.ConsumptionBased = types.BoolValue(cliRes.ConsumptionBased)
	data.ContractStatus = types.StringValue(cliRes.ContractStatus)
	data.CostObjectId = stringNullIfEmpty(cliRes.CostObjectId)
	data.CostObjectType = stringNullIfEmpty(cliRes.CostObjectType)
	data.CreatedDate = timeToValue(cliRes.CreatedDate.Time())

	data.CrmCustomerId = stringNullIfEmpty(cliRes.CrmCustomerId)
	data.CrmTenantId = stringNullIfEmpty(cliRes.CrmTenantId)

	data.Description = types.StringValue(cliRes.Description)
	data.DisplayName = types.StringValue(cliRes.DisplayName)
	data.ExpiryDate = timeToValue(cliRes.ExpiryDate.Time())
	data.GeoAccess = types.StringValue(cliRes.GeoAccess)
	data.LicenseType = types.StringValue(cliRes.LicenseType)
	data.LastModified = timeToValue(cliRes.ModifiedDate.Time())
	data.State = types.StringValue(cliRes.EntityState)
	data.Origin = types.StringValue(cliRes.Origin)
	data.RenewalDate = timeToValue(cliRes.RenewalDate.Time())

	data.ServiceId = stringNullIfEmpty(cliRes.ServiceId)

	data.Subdomain = types.StringValue(cliRes.Subdomain)
	data.Usage = stringNullIfEmpty(cliRes.UseFor)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
