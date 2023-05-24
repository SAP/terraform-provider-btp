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

__Tip__
You must be assigned to the global account admin or viewer role.

__Further documentation__
https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/8ed4a705efa0431b910056c0acdbf377.html`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique ID of the global account.",
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
				MarkdownDescription: "Whether the customer of the global account pays only for services that they actually use (consumption-based) or pay for subscribed services at a fixed cost irrespective of consumption (subscription-based). * <b>TRUE:</b> Consumption-based commercial model. * <b>FALSE:</b> Subscription-based commercial model.",
				Computed:            true,
			},
			"contract_status": schema.StringAttribute{
				MarkdownDescription: "The status of the customer contract and its associated root global account. * <b>ACTIVE:</b> The customer contract and its associated global account is currently active. * <b>PENDING_TERMINATION:</b> A termination process has been triggered for a customer contract (the customer contract has expired, or a customer has given notification that they wish to terminate their contract), and the global account is currently in the validation period. The customer can still access their global account until the end of the validation period. * <b>SUSPENDED:</b> For enterprise accounts, specifies that the customer's global account is currently in the grace period of the termination process. Access to the global account by the customer is blocked. No data is deleted until the deletion date is reached at the end of the grace period. For trial accounts, specifies that the account is suspended, and the account owner has not yet extended the trial period.",
				Computed:            true,
			},
			"costobject_id": schema.StringAttribute{
				MarkdownDescription: "The number or code of the cost center, internal order, or Work Breakdown Structure element that is charged for the creation and usage of the global account. The type of the cost object must be configured in costObjectType.",
				Computed:            true,
			},
			"costobject_type": schema.StringAttribute{
				MarkdownDescription: "The type of accounting assignment object that is associated with the global account owner and used to charge for the creation and usage of the global account. Support types: COST_CENTER, INTERNAL_ORDER, WBS_ELEMENT. The number or code of the specified cost object is defined in costObjectId. For a cost object of type 'cost center', the value is also configured in costCenter for backward compatibility purposes.",
				Computed:            true,
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
				MarkdownDescription: "The geographic locations from where the global account can be accessed. * <b>STANDARD:</b> The global account can be accessed from any geographic location. * <b>EU_ACCESS:</b> The global account can be accessed only within locations in the EU.",
				Computed:            true,
			},
			"license_type": schema.StringAttribute{
				MarkdownDescription: "The type of license for the global account. The license type affects the scope of functions of the account. * <b>DEVELOPER:</b> For internal developer global accounts on Staging or Canary landscapes. * <b>CUSTOMER:</b> For customer global accounts. * <b>PARTNER:</b> For partner global accounts. * <b>INTERNAL_DEV:</b> For internal global accounts on the Dev landscape. * <b>INTERNAL_PROD:</b> For internal global accounts on the Live landscape. * <b>TRIAL:</b> For customer trial accounts.",
				Computed:            true,
			},
			"origin": schema.StringAttribute{
				MarkdownDescription: "The origin of the account. * <b>ORDER:</b> Created by the Order Processing API or Submit Order wizard. * <b>OPERATOR:</b> Created by the Global Account wizard. * <b>REGION_SETUP:</b> Created automatically as part of the region setup.",
				Computed:            true,
			},
			"service_id": schema.StringAttribute{
				MarkdownDescription: "For internal accounts, the service for which the global account was created.",
				Computed:            true,
			},
			"subdomain": schema.StringAttribute{
				MarkdownDescription: "Relevant only for entities that require authorization (e.g. global account). The subdomain that becomes part of the path used to access the authorization tenant of the global account. Unique within the defined region.",
				Computed:            true,
			},
			"usage": schema.StringAttribute{
				MarkdownDescription: "For internal accounts, the intended purpose of the global account. Possible purposes: * <b>Development:</b> For development of a service. * <b>Testing:</b> For testing development. * <b>Demo:</b> For creating demos. * <b>Production:</b> For delivering a service in a production landscape.",
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The current state of the global account. * <b>STARTED:</b> CRUD operation on an entity has started. * <b>CREATING:</b> Creating entity operation is in progress. * <b>UPDATING:</b> Updating entity operation is in progress. * <b>MOVING:</b> Moving entity operation is in progress. * <b>PROCESSING:</b> A series of operations related to the entity is in progress. * <b>DELETING:</b> Deleting entity operation is in progress. * <b>OK:</b> The CRUD operation or series of operations completed successfully. * <b>PENDING REVIEW:</b> The processing operation has been stopped for reviewing and can be restarted by the operator. * <b>CANCELLED:</b> The operation or processing was canceled by the operator. * <b>CREATION_FAILED:</b> The creation operation failed, and the entity was not created or was created but cannot be used. * <b>UPDATE_FAILED:</b> The update operation failed, and the entity was not updated. * <b>PROCESSING_FAILED:</b> The processing operations failed. * <b>DELETION_FAILED:</b> The delete operation failed, and the entity was not deleted. * <b>MOVE_FAILED:</b> Entity could not be moved to a different location. * <b>MIGRATING:</b> Migrating entity from NEO to CF.",
				Computed:            true,
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
