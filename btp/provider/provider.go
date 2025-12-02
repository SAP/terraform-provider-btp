package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/version"
)

const userPasswordFlow = "userPasswordFlow"
const x509Flow = "x509Flow"
const idTokenFlow = "idTokenFlow"
const ssoFlow = "ssoFlow"
const assertionFlow = "assertionFlow"
const errorMessagePostfixWithEnv = "If either is already set, ensure the value is not empty."
const errorMessagePostfixWithoutEnv = "If it is already set, ensure the value is not empty."

func New() provider.Provider {
	return NewWithClient(http.DefaultClient)
}

func NewWithClient(httpClient *http.Client) provider.Provider {
	return &btpcliProvider{httpClient: httpClient}
}

func NewWithFunctions() provider.ProviderWithFunctions {
	return NewWithFunctionsAndClient(http.DefaultClient)
}

func NewWithFunctionsAndClient(httpClient *http.Client) provider.ProviderWithFunctions {
	return &btpcliProvider{httpClient: httpClient}
}

type btpcliProvider struct {
	httpClient          *http.Client
	betaFeaturesEnabled bool
}

func (p *btpcliProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The Terraform provider for SAP BTP enables you to automate the provisioning, management, and configuration of resources on [SAP Business Technology Platform](https://account.hana.ondemand.com/). By leveraging this provider, you can simplify and streamline the deployment and maintenance of BTP services and applications.`,
		Attributes: map[string]schema.Attribute{
			"cli_server_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the BTP CLI server (e.g. `https://cli.btp.cloud.sap`).",
				Optional:            true, // TODO validate URL
			},
			"globalaccount": schema.StringAttribute{
				MarkdownDescription: "The subdomain of the global account in which you want to manage resources. To be found in the cockpit, in the global account view.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Your user name, usually an e-mail address. This can also be sourced from the `BTP_USERNAME` environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Your password. Note that two-factor authentication is not supported. This can also be sourced from the `BTP_PASSWORD` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"idtoken": schema.StringAttribute{
				MarkdownDescription: "A valid id token. To be provided instead of 'username' and 'password'. This can also be sourced from the `BTP_IDTOKEN` environment variable. (SAP-internal usage only)",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("username"), path.MatchRoot("password"), path.MatchRoot("idp"), path.MatchRoot("tls_idp_url"), path.MatchRoot("tls_client_key"), path.MatchRoot("tls_client_certificate")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"assertion": schema.StringAttribute{
				MarkdownDescription: "A valid assertion JWT token. To be provided instead of 'username' and 'password'. This can also be sourced from the `BTP_ASSERTION` environment variable. This authentication method is only supported when using a custom Identity Provider (IdP).",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("username"), path.MatchRoot("password"), path.MatchRoot("tls_idp_url"), path.MatchRoot("tls_client_key"), path.MatchRoot("tls_client_certificate"), path.MatchRoot("idtoken")),
					stringvalidator.AlsoRequires(path.MatchRoot("idp")),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"idp": schema.StringAttribute{
				MarkdownDescription: "The identity provider to be used for authentication (only required for custom idp).",
				Optional:            true,
			},
			"tls_idp_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the identity provider to be used for authentication (only required for x509 auth).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("password"), path.MatchRoot("idtoken")),
					stringvalidator.AlsoRequires(path.MatchRoot("tls_client_key"), path.MatchRoot("tls_client_certificate"), path.MatchRoot("idp")),
				},
			},
			"tls_client_key": schema.StringAttribute{
				MarkdownDescription: "PEM encoded private key (only required for x509 auth).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("password"), path.MatchRoot("idtoken")),
					stringvalidator.AlsoRequires(path.MatchRoot("tls_idp_url"), path.MatchRoot("tls_client_certificate"), path.MatchRoot("idp")),
				},
			},
			"tls_client_certificate": schema.StringAttribute{
				MarkdownDescription: "PEM encoded certificate (only required for x509 auth).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("password"), path.MatchRoot("idtoken")),
					stringvalidator.AlsoRequires(path.MatchRoot("tls_idp_url"), path.MatchRoot("tls_client_key"), path.MatchRoot("idp")),
				},
			},
		},
	}
}

type providerData struct {
	CLIServerURL         types.String `tfsdk:"cli_server_url"`
	GlobalAccount        types.String `tfsdk:"globalaccount"`
	Username             types.String `tfsdk:"username"`
	Password             types.String `tfsdk:"password"`
	Assertion            types.String `tfsdk:"assertion"`
	IdToken              types.String `tfsdk:"idtoken"`
	IdentityProvider     types.String `tfsdk:"idp"`
	IdentityProviderURL  types.String `tfsdk:"tls_idp_url"`
	TLSClientKey         types.String `tfsdk:"tls_client_key"`
	TLSClientCertificate types.String `tfsdk:"tls_client_certificate"`
}

// Metadata returns the provider type name.
func (p *btpcliProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "btp"
}

func (p *btpcliProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	const unableToCreateClient = "unableToCreateClient"

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
		resp.Diagnostics.AddError(unableToCreateClient, fmt.Sprintf("%s", err))
		return
	}

	client := btpcli.NewClientFacade(btpcli.NewV2ClientWithHttpClient(p.httpClient, u, nil))
	btpUserAgent := os.Getenv("BTP_APPEND_USER_AGENT")

	if len(strings.TrimSpace(btpUserAgent)) == 0 {
		client.UserAgent = fmt.Sprintf("Terraform/%s terraform-provider-btp/%s", req.TerraformVersion, version.ProviderVersion)
	} else {
		client.UserAgent = fmt.Sprintf("Terraform/%s terraform-provider-btp/%s custom-user-agent/%s", req.TerraformVersion, version.ProviderVersion, btpUserAgent)
	}

	ssoLogin := false
	enableSSO := os.Getenv("BTP_ENABLE_SSO")
	if len(strings.TrimSpace(enableSSO)) != 0 {
		ssoLogin, err = strconv.ParseBool(enableSSO)
		if err != nil {
			resp.Diagnostics.AddError("unable to convert sso value", fmt.Sprintf("%s", err))
			return
		}
	}

	// User may provide an idp to the provider
	var idp string
	if config.IdentityProvider.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as identity provider")
		return
	}

	if config.IdentityProvider.IsNull() {
		idp = os.Getenv("BTP_IDP")
	} else {
		idp = config.IdentityProvider.ValueString()
	}

	// User may provide an id token to the provider instead of username and password (see below)
	var idToken string
	if config.IdToken.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as id token")
		return
	}
	//BTP_ASSERTION
	if config.Assertion.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as assertion")
		return
	}
	var assertion string
	if config.Assertion.IsNull() {
		assertion = os.Getenv("BTP_ASSERTION")
	} else {
		assertion = config.Assertion.ValueString()
	}

	if config.IdToken.IsNull() {
		idToken = os.Getenv("BTP_IDTOKEN")
	} else {
		idToken = config.IdToken.ValueString()
	}

	// User must provide a username to the provider unless an id token is given
	var username string
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as username")
		return
	}

	if config.Username.IsNull() {
		username = os.Getenv("BTP_USERNAME")
	} else {
		username = config.Username.ValueString()
	}

	// User must provide a password to the provider unless an id token is given
	var password string
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as password")
		return
	}

	if config.Password.IsNull() {
		password = os.Getenv("BTP_PASSWORD")
	} else {
		password = config.Password.ValueString()
	}

	//Check for conflicts between the different auth flows
	//This can happen if the user proivdes the values via ENV variables as the schema validation will not catch this
	if len(idToken) > 0 && (len(username) > 0 || len(password) > 0) {
		resp.Diagnostics.AddError(unableToCreateClient, "Cannot provide both id token and username/password")
		return
	}

	//Determine and execute the login flow depending on the provided parameters
	switch authFlow := determineAuthFlow(config, idToken, ssoLogin, assertion); authFlow {
	case ssoFlow:
		if _, err = client.BrowserLogin(ctx, btpcli.NewBrowserLoginRequest(idp, config.GlobalAccount.ValueString())); err != nil {
			resp.Diagnostics.AddError(unableToCreateClient, fmt.Sprintf("%s", err))
		}
	case userPasswordFlow:
		validateUserPasswordFlow(username, password, resp)

		if resp.Diagnostics.HasError() {
			return
		}

		if _, err = client.Login(ctx, btpcli.NewLoginRequestWithCustomIDP(idp, config.GlobalAccount.ValueString(), username, password)); err != nil {
			resp.Diagnostics.AddError(unableToCreateClient, fmt.Sprintf("%s", err))
		}
	case x509Flow:

		validateX509Flow(username, config.IdentityProviderURL.ValueString(), config.TLSClientKey.ValueString(), config.TLSClientCertificate.ValueString(), resp)

		if resp.Diagnostics.HasError() {
			return
		}

		passcodeLoginReq := &btpcli.PasscodeLoginRequest{
			GlobalAccountSubdomain: config.GlobalAccount.ValueString(),
			IdentityProvider:       idp,
			IdentityProviderURL:    config.IdentityProviderURL.ValueString(),
			Username:               username,
			PEMEncodedPrivateKey:   config.TLSClientKey.ValueString(),
			PEMEncodedCertificate:  config.TLSClientCertificate.ValueString(),
		}

		if _, err = client.PasscodeLogin(ctx, passcodeLoginReq); err != nil {
			resp.Diagnostics.AddError(unableToCreateClient, fmt.Sprintf("%s", err))
			return
		}
	case assertionFlow:
		if _, err = client.Login(ctx, btpcli.NewLoginRequestWithAssertion(idp, config.GlobalAccount.ValueString(), assertion)); err != nil {
			resp.Diagnostics.AddError(unableToCreateClient, fmt.Sprintf("%s", err))
			return
		}

	case idTokenFlow:
		// SAP Internal usage only
		if _, err = client.IdTokenLogin(ctx, btpcli.NewIdTokenLoginRequest(config.GlobalAccount.ValueString(), idToken)); err != nil {
			resp.Diagnostics.AddError(unableToCreateClient, fmt.Sprintf("%s", err))
		}
	default:
		// No valid login flow
		resp.Diagnostics.AddError(unableToCreateClient, "No valid login flow found. Please provide either username and password, or an id token, or a client certificate and key.")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// Resources - Defines provider resources
func (p *btpcliProvider) Resources(ctx context.Context) []func() resource.Resource {
	betaResources := []func() resource.Resource{
		//Beta resources should be excluded from sonar scan.
		//If you add them to production code, remove them from sonar exclusion list
	}

	if !p.betaFeaturesEnabled {
		betaResources = nil
	}

	return append([]func() resource.Resource{
		newDirectoryApiCredentialResource,
		newDirectoryResource,
		newDirectoryEntitlementResource,
		newDirectoryRoleCollectionAssignmentResource,
		newDirectoryRoleCollectionResource,
		newGlobalaccountApiCredentialResource,
		newGlobalaccountResourceProviderResource,
		newGlobalaccountRoleCollectionAssignmentResource,
		newGlobalaccountRoleCollectionResource,
		newGlobalaccountSecuritySettingsResource,
		newGlobalaccountTrustConfigurationResource,
		newSubaccountApiCredentialResource,
		newSubaccountEntitlementResource,
		newSubaccountEnvironmentInstanceResource,
		newSubaccountResource,
		newSubaccountRoleCollectionAssignmentResource,
		newSubaccountRoleCollectionResource,
		newSubaccountSecuritySettingsResource,
		newSubaccountServiceBindingResource,
		newSubaccountServiceBrokerResource,
		newSubaccountServiceInstanceResource,
		newSubaccountSubscriptionResource,
		newSubaccountTrustConfigurationResource,
		newDirectoryRoleResource,
		newGlobalaccountRoleResource,
		newSubaccountRoleResource,
	}, betaResources...)
}

// DataSources - Defines provider data sources
func (p *btpcliProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	betaDataSources := []func() datasource.DataSource{
		//Beta data sources should be excluded from sonar scan.
		//If you add them to production code, remove them from sonar exclusion list
		newDirectoryAppDataSource,
		newDirectoryAppsDataSource,
		newGlobalaccountAppDataSource,
		newGlobalaccountAppsDataSource,
		newGlobalaccountResourceProviderDataSource,
		newGlobalaccountResourceProvidersDataSource,
		newSubaccountServicePlatformDataSource,
		newSubaccountServicePlatformsDataSource,
	}

	if !p.betaFeaturesEnabled {
		betaDataSources = nil
	}

	return append([]func() datasource.DataSource{
		newDirectoryDataSource,
		newDirectoriesDataSource,
		newDirectoryEntitlementDataSource,
		newDirectoryEntitlementsDataSource,
		newDirectoryLabelsDataSource,
		newDirectoryRoleCollectionDataSource,
		newDirectoryRoleCollectionsDataSource,
		newDirectoryRoleDataSource,
		newDirectoryRolesDataSource,
		newDirectoryUserDataSource,
		newDirectoryUsersDataSource,
		newGlobalaccountDataSource,
		newGlobalaccountWithHierarchyDataSource,
		newGlobalaccountEntitlementsDataSource,
		newGlobalaccountRoleCollectionDataSource,
		newGlobalaccountRoleCollectionsDataSource,
		newGlobalaccountRoleDataSource,
		newGlobalaccountRolesDataSource,
		newGlobalaccountSecuritySettingsDataSource,
		newGlobalaccountTrustConfigurationDataSource,
		newGlobalaccountTrustConfigurationsDataSource,
		newGlobalaccountUserDataSource,
		newGlobalaccountUsersDataSource,
		newRegionsDataSource,
		newSubaccountAppDataSource,
		newSubaccountAppsDataSource,
		newSubaccountDataSource,
		newSubaccountEntitlementDataSource,
		newSubaccountEntitlementsDataSource,
		newSubaccountEnvironmentInstanceDataSource,
		newSubaccountEnvironmentInstancesDataSource,
		newSubaccountEnvironmentsDataSource,
		newSubaccountLabelsDataSource,
		newSubaccountRoleCollectionDataSource,
		newSubaccountRoleCollectionsDataSource,
		newSubaccountRoleDataSource,
		newSubaccountRolesDataSource,
		newSubaccountSecuritySettingsDataSource,
		newSubaccountServiceBindingDataSource,
		newSubaccountServiceBindingsDataSource,
		newSubaccountServiceBrokerDataSource,
		newSubaccountServiceBrokersDataSource,
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
		newSubaccountDestinationTrustDataSource,
		newSubaccountDestinationFragmentDataSource,
		newSubaccountDestinationFragmentsDataSource,
	}, betaDataSources...)
}

func (p *btpcliProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewExtractCfApiUrlFunction,
		NewExtractCfOrgIdFunction,
		NewExtractKymaApiServerUrlFunction,
	}
}

func determineAuthFlow(config providerData, idToken string, ssoLogin bool, assertion string) string {
	if ssoLogin {
		return ssoFlow
	} else if len(idToken) > 0 {
		return idTokenFlow
	} else if !config.TLSClientKey.IsNull() {
		return x509Flow
	} else if len(assertion) > 0 {
		return assertionFlow
	} else {
		return userPasswordFlow
	}
}

func validateUserPasswordFlow(userName string, password string, resp *provider.ConfigureResponse) {
	if len(userName) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Username",
			"The provider cannot create the Terraform BTP client as there is a missing or empty value for the username. "+
				"Set the username value in the configuration or use the BTP_USERNAME environment variable. "+
				errorMessagePostfixWithEnv,
		)
	}

	if len(password) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Password",
			"The provider cannot create the Terraform BTP client as there is a missing or empty value for the password. "+
				"Set the password value in the configuration or use the BTP_PASSWORD environment variable. "+
				errorMessagePostfixWithEnv,
		)

	}

}

func validateX509Flow(userName string, identityProviderUrl string, tlsClientKey string, tlsClientCertificate string, resp *provider.ConfigureResponse) {
	if len(userName) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Username",
			"The provider cannot create the Terraform BTP client as there is a missing or empty value for the username. "+
				"Set the username value in the configuration or use the BTP_USERNAME environment variable. "+
				errorMessagePostfixWithEnv,
		)
	}

	if len(identityProviderUrl) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("tls_idp_url"),
			"Missing TLS IDP URL (only required for x509 auth)",
			"The provider cannot create the Terraform BTP client as there is a missing or empty value for the tls_idp_url (only required for x509 auth). "+
				"Set the tls_idp_url value in the configuration. "+
				errorMessagePostfixWithoutEnv,
		)
	}

	if len(tlsClientKey) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("tls_client_key"),
			"Missing PEM Encoded Private Key",
			"The provider cannot create the Terraform BTP client as there is a missing or empty value for the tls_client_key (PEM encoded private key). "+
				"Set the tls_client_key value in the configuration. "+
				errorMessagePostfixWithoutEnv,
		)
	}

	if len(tlsClientCertificate) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("tls_client_certificate"),
			"Missing PEM Encoded Certificate",
			"The provider cannot create the Terraform BTP client as there is a missing or empty value for the tls_client_certificate (PEM encoded certificate). "+
				"Set the tls_client_certificate value in the configuration. "+
				errorMessagePostfixWithoutEnv,
		)
	}
}
