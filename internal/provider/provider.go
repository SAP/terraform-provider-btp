package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
)

// New .
func New() provider.Provider {
	return NewWithClient(http.DefaultClient)
}

func NewWithClient(httpClient *http.Client) provider.Provider {
	return &btpcliProvider{httpClient: httpClient}
}

type btpcliProvider struct {
	httpClient          *http.Client
	betaFeaturesEnabled bool
}

// GetSchema
func (p *btpcliProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The Terraform provider for SAP BTP enables you to automate the provisioning, management, and configuration of resources on [SAP Business Technology Platform](https://account.hana.ondemand.com/). By leveraging this provider, you can simplify and streamline the deployment and maintenance of BTP services and applications.`,
		Attributes: map[string]schema.Attribute{
			"cli_server_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the BTP CLI server (e.g. `https://cpcli.cf.eu10.hana.ondemand.com`).",
				Optional:            true, // TODO validate URL
			},
			"globalaccount": schema.StringAttribute{
				MarkdownDescription: "The subdomain of the global account you want to log in to. To be found in the cockpit, in the global account view.",
				Required:            true, // TODO validate UUID
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Your user name, usually an e-mail address.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Your password. Note that if two-factor authentication is enabled, concatenate your password, followed by the passcode, in a single string.",
				Optional:            true,
				Sensitive:           true,
			},
			"auth": schema.StringAttribute{
				MarkdownDescription: "Select the authentication method of your choice. Can be either `sso` or `password` (default).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("sso", "password"),
				},
			},
		},
	}
}

// Provider schema struct
type providerData struct {
	CLIServerURL  types.String `tfsdk:"cli_server_url"`
	GlobalAccount types.String `tfsdk:"globalaccount"`
	Username      types.String `tfsdk:"username"`
	Password      types.String `tfsdk:"password"`
	Auth          types.String `tfsdk:"auth"`
}

// Metadata returns the provider type name.
func (p *btpcliProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "btp"
}

func (p *btpcliProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	selectedCLIServerURL := btpcli.DefaultServerURL

	if !config.CLIServerURL.IsNull() {
		selectedCLIServerURL = config.CLIServerURL.ValueString()
	}

	u, err := url.Parse(selectedCLIServerURL) // TODO move to NewV2Client

	if err != nil {
		resp.Diagnostics.AddError("Unable to create Client", fmt.Sprintf("%s", err))
		return
	}

	client := btpcli.NewClientFacade(btpcli.NewV2ClientWithHttpClient(p.httpClient, u))

	// User must provide a username to the provider
	var username string
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddWarning("Unable to create client", "Cannot use unknown value as client_certificate")
		return
	}

	if config.Username.IsNull() {
		username = os.Getenv("BTP_USERNAME")
	} else {
		username = config.Username.ValueString()
	}

	// User must provide a password to the provider
	var password string
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddWarning("Unable to create client", "Cannot use unknown value as password")
		return
	}

	if config.Password.IsNull() {
		password = os.Getenv("BTP_PASSWORD")
	} else {
		password = config.Password.ValueString()
	}

	if len(username) == 0 || len(password) == 0 {
		resp.Diagnostics.AddError("Unable to create Client", "globalaccount, username and password must be given.")
		return
	}

	if _, err = client.Login(ctx, btpcli.NewLoginRequest(config.GlobalAccount.ValueString(), username, password)); err != nil {
		resp.Diagnostics.AddError("Unable to create Client", fmt.Sprintf("%s", err))
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources - Defines provider resources
func (p *btpcliProvider) Resources(ctx context.Context) []func() resource.Resource {
	betaResources := []func() resource.Resource{
		newDirectoryRoleResource,
	}

	if !p.betaFeaturesEnabled {
		betaResources = nil
	}

	return append([]func() resource.Resource{
		newDirectoryResource,
		newDirectoryRoleCollectionResource,
		newDirectoryRoleCollectionAssignmentResource,
		newGlobalaccountResourceProviderResource,
		newGlobalaccountRoleCollectionResource,
		newGlobalaccountRoleCollectionAssignmentResource,
		newGlobalaccountRoleResource,
		newGlobalaccountTrustConfigurationResource,
		newSubaccountEntitlementResource,
		newSubaccountEnvironmentInstanceResource,
		newSubaccountResource,
		newSubaccountRoleCollectionResource,
		newSubaccountRoleCollectionAssignmentResource,
		newSubaccountRoleResource,
		newSubaccountServiceBindingResource,
		newSubaccountServiceInstanceResource,
		newSubaccountSubscriptionResource,
		newSubaccountTrustConfigurationResource,
	}, betaResources...)
}

// DataSources - Defines provider data sources
func (p *btpcliProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	betaDataSources := []func() datasource.DataSource{
		newDirectoryAppDataSource,
		newDirectoryAppsDataSource,
		newGlobalaccountAppDataSource,
		newGlobalaccountAppsDataSource,
		newSubaccountAppDataSource,
		newSubaccountAppsDataSource,
		newSubaccountServiceBrokerDataSource,
		newSubaccountServiceBrokersDataSource,
		newSubaccountServicePlatformDataSource,
		newSubaccountServicePlatformsDataSource,
	}

	if !p.betaFeaturesEnabled {
		betaDataSources = nil
	}

	return append([]func() datasource.DataSource{
		newDirectoryDataSource,
		newDirectoryEntitlementsDataSource,
		newDirectoryLabelsDataSource,
		newDirectoryRoleCollectionDataSource,
		newDirectoryRoleCollectionsDataSource,
		newDirectoryRoleDataSource,
		newDirectoryRolesDataSource,
		newDirectoryUserDataSource,
		newDirectoryUsersDataSource,
		newGlobalaccountDataSource,
		newGlobalaccountEntitlementsDataSource,
		newGlobalaccountResourceProviderDataSource,
		newGlobalaccountResourceProvidersDataSource,
		newGlobalaccountRoleCollectionDataSource,
		newGlobalaccountRoleCollectionsDataSource,
		newGlobalaccountRoleDataSource,
		newGlobalaccountRolesDataSource,
		newGlobalaccountTrustConfigurationDataSource,
		newGlobalaccountTrustConfigurationsDataSource,
		newGlobalaccountUserDataSource,
		newGlobalaccountUsersDataSource,
		newRegionsDataSource,
		newSubaccountDataSource,
		newSubaccountEntitlementsDataSource,
		newSubaccountEnvironmentInstanceDataSource,
		newSubaccountEnvironmentInstancesDataSource,
		newSubaccountEnvironmentsDataSource,
		newSubaccountLabelsDataSource,
		newSubaccountRoleCollectionDataSource,
		newSubaccountRoleCollectionsDataSource,
		newSubaccountRoleDataSource,
		newSubaccountRolesDataSource,
		newSubaccountServiceBindingDataSource,
		newSubaccountServiceBindingsDataSource,
		newSubaccountServiceInstanceDataSource,
		newSubaccountServiceInstancesDataSource,
		newSubaccountServiceOfferingDataSource,
		newSubaccountServiceOfferingsDataSource,
		newSubaccountServicePlanDataSource,
		newSubaccountServicePlansDataSource,
		newSubaccountSubscriptionDataSource,
		newSubaccountSubscriptionsDataSource,
		newSubaccountTrustConfigurationDataSource,
		newSubaccountTrustConfigurationsDataSource,
		newSubaccountUserDataSource,
		newSubaccountUsersDataSource,
		newSubaccountsDataSource,
		newWhoamiDataSource,
	}, betaDataSources...)
}
