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
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

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
const btpCliSessionFlow = "btpCliSessionFlow"
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

func NewWithActions() provider.ProviderWithActions {
	return NewWithActionsAndClient(http.DefaultClient)
}

func NewWithActionsAndClient(httpClient *http.Client) provider.ProviderWithActions {
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
		if ssoLogin {
			resp.Diagnostics.AddWarning("deprecated authentication flow", "Do not use SSO login, use the login via the btp CLI session instead.")
		}
	}

	btpCliSessionLogin := false
	btpCliCustomConfigPath := ""
	enableBTPCliSessionLogin := os.Getenv("USE_BTPCLI_SESSION")
	if len(strings.TrimSpace(enableBTPCliSessionLogin)) != 0 {
		btpCliSessionLogin, err = strconv.ParseBool(enableBTPCliSessionLogin)
		if err != nil {
			resp.Diagnostics.AddError("unable to convert btp cli session login value", fmt.Sprintf("%s", err))
			return
		}
		btpCliCustomConfigPath = os.Getenv("BTPCLI_CONFIG_PATH")
	}

	// Explicit provider attributes take precedence over the corresponding ENV
	// variable. When both are set, resolveWithEnv emits a warning naming the
	// attribute that wins and the env var that is ignored.
	if config.IdentityProvider.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as identity provider")
		return
	}
	if config.IdToken.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as id token")
		return
	}
	if config.Assertion.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as assertion")
		return
	}
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as username")
		return
	}
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddWarning(unableToCreateClient, "Cannot use unknown value as password")
		return
	}

	idp, _ := resolveWithEnv(config.IdentityProvider, "BTP_IDP", "idp", resp)
	assertion, assertionExplicit := resolveWithEnv(config.Assertion, "BTP_ASSERTION", "assertion", resp)
	idToken, idTokenExplicit := resolveWithEnv(config.IdToken, "BTP_IDTOKEN", "idtoken", resp)
	username, usernameExplicit := resolveWithEnv(config.Username, "BTP_USERNAME", "username", resp)
	password, passwordExplicit := resolveWithEnv(config.Password, "BTP_PASSWORD", "password", resp)

	// Cross-flow conflict: when the user picked a flow via any explicit attribute,
	// values sourced from env vars that belong to *other* flows must be dropped
	// with a warning (schema ConflictsWith only covers explicit-vs-explicit).
	username, password, idToken, assertion = dropCrossFlowEnvValues(
		config,
		username, usernameExplicit,
		password, passwordExplicit,
		idToken, idTokenExplicit,
		assertion, assertionExplicit,
		resp,
	)

	//Check for conflicts between the different auth flows
	//This can happen if the user provides the values via ENV variables as the schema validation will not catch this
	if len(idToken) > 0 && (len(username) > 0 || len(password) > 0) {
		resp.Diagnostics.AddError(unableToCreateClient, "Cannot provide both id token and username/password")
		return
	}

	if btpCliSessionLogin && (len(username) > 0 || len(password) > 0) {
		resp.Diagnostics.AddError(unableToCreateClient, "Cannot provide both BTP CLI session login and username/password")
		return
	}

	// Log resolved provider configuration at DEBUG level to aid troubleshooting
	// and quickly identify basic misconfigurations (e.g. Global Account, IdP,
	// or CLI server URL) during provider initialization.
	tflog.Debug(ctx, "Initializing SAP BTP provider with resolved configuration as:", map[string]any{
		"global_account":     config.GlobalAccount.ValueString(),
		"idp":                idp,
		"cli_server_url":     selectedCLIServerURL,
		"use_btpcli_session": btpCliSessionLogin,
	})

	//Determine and execute the login flow depending on the provided parameters
	switch authFlow := determineAuthFlow(config, idToken, ssoLogin, assertion, btpCliSessionLogin); authFlow {
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

	case btpCliSessionFlow:
		if _, err = client.BtpCliSessionLogin(ctx, btpcli.NewBtpCliSessionLoginRequest(config.GlobalAccount.ValueString(), btpCliCustomConfigPath, idp)); err != nil {
			resp.Diagnostics.AddError(unableToCreateClient, fmt.Sprintf("%s", err))
			return
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
	resp.ListResourceData = client
	resp.ActionData = client
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
		newSubaccountDestinationFragmentResource,
		newSubaccountEntitlementResource,
		newSubaccountEnvironmentInstanceResource,
		newSubaccountResource,
		newSubaccountRoleCollectionAssignmentResource,
		newSubaccountRoleCollectionResource,
		newSubaccountRoleCollectionRoleResource,
		newSubaccountRoleCollectionBaseResource,
		newSubaccountSecuritySettingsResource,
		newSubaccountServiceBindingResource,
		newSubaccountServiceBrokerResource,
		newSubaccountServiceInstanceResource,
		newSubaccountSubscriptionResource,
		newSubaccountTrustConfigurationResource,
		newDirectoryRoleResource,
		newGlobalaccountRoleResource,
		newSubaccountRoleResource,
		newSubaccountDestinationResource,
		newSubaccountDestinationCertificateResource,
		newSubaccountDestinationGenericResource,
		newDisasterRecoverySubaccountPairResource,
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
		newGlobalaccountEntitlementsWithDcDataSource,
		newGlobalaccountEntitlementWithDcDataSource,
		newGlobalaccountRoleCollectionDataSource,
		newGlobalaccountRoleCollectionsDataSource,
		newGlobalaccountRoleDataSource,
		newGlobalaccountRolesDataSource,
		newGlobalaccountSecuritySettingsDataSource,
		newGlobalaccountSecurityIdentityProviderDataSource,
		newGlobalaccountSecurityIdentityProvidersDataSource,
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
		newSubaccountRoleCollectionRoleDataSource,
		newSubaccountRoleCollectionRolesDataSource,
		newSubaccountRoleCollectionBaseDataSource,
		newSubaccountRoleCollectionBasesDataSource,
		newSubaccountRoleDataSource,
		newSubaccountRolesDataSource,
		newSubaccountSecuritySettingsDataSource,
		newSubaccountSecurityIdentityProviderDataSource,
		newSubaccountSecurityIdentityProvidersDataSource,
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
		newSubaccountDestinationCertificateDataSource,
		newSubaccountDestinationCertificatesDataSource,
		newSubaccountDestinationDataSource,
		newSubaccountDestinationsDataSource,
		newSubaccountDestinationTrustDataSource,
		newSubaccountDestinationFragmentDataSource,
		newSubaccountDestinationFragmentsDataSource,
		newSubaccountDestinationsGenericDataSource,
		newSubaccountDestinationsNamesDataSource,
		newSubaccountDestinationGenericDataSource,
		newDisasterRecoverySubaccountPairDataSource,
	}, betaDataSources...)
}

func (p *btpcliProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewExtractCfApiUrlFunction,
		NewExtractCfOrgIdFunction,
		NewExtractKymaApiServerUrlFunction,
		NewExtractKymaKubeconfigUrlFunction,
		NewDownloadKymaKubeconfigFunction,
	}
}

// ListResources defines the ListResources implemented in the provider.
func (p *btpcliProvider) ListResources(_ context.Context) []func() list.ListResource {
	return []func() list.ListResource{
		NewGlobalaccountRoleListResource,
		NewGlobalaccountResourceProviderListResource,
		NewGlobalaccountRoleCollectionListResource,
		NewDirectoryEntitlementListResource,
		NewDirectoryRoleCollectionListResource,
		NewSubaccountServiceBrokerListResource,
		NewSubaccountServiceInstanceListResource,
		NewSubaccountEnvironmentInstanceListResource,
		NewSubaccountListResource,
		NewGlobalaccountTrustConfigurationListResource,
		NewSubaccountTrustConfigurationListResource,
		NewSubaccountServiceBindingListResource,
		NewSubaccountSecuritySettingsListResource,
		NewGlobalaccountSecuritySettingsListResource,
		NewDirectoryListResource,
		NewDirectoryRoleListResource,
		NewSubaccountDestinationGenericListResource,
		NewSubaccountSubscriptionListResource,
		NewSubaccountRoleCollectionListResource,
		NewSubaccountRoleListResource,
		NewSubaccountDestinationFragmentListResource,
		NewSubaccountEntitlementListResource,
		NewSubaccountRoleCollectionRoleListResource,
		NewSubaccountRoleCollectionBaseListResource,
	}
}

// ActionsResources defines the Terraform Actions implemented in the provider.
func (p *btpcliProvider) Actions(_ context.Context) []func() action.Action {
	return []func() action.Action{
		NewRestoreSubaccountAction,
		NewAddMeAsSubaccountAdminAction,
	}
}

// resolveWithEnv returns the value and whether it came from the explicit
// provider attribute (true) or from the env var / defaulted (false).
// Explicit attribute wins over env; a warning is emitted when both are set.
// ponytail: null == "not provided"; explicit empty string still falls through
// to env with no warning (matches prior behavior).
func resolveWithEnv(cfg types.String, envName, attrName string, resp *provider.ConfigureResponse) (string, bool) {
	envVal := os.Getenv(envName)
	if cfg.IsNull() {
		return envVal, false
	}
	explicit := cfg.ValueString()
	if len(explicit) > 0 && len(strings.TrimSpace(envVal)) > 0 {
		resp.Diagnostics.AddWarning(
			"Conflicting authentication configuration",
			fmt.Sprintf("Both the provider attribute %q and the environment variable %q are set. The explicit provider attribute takes precedence; %q is ignored.", attrName, envName, envName),
		)
	}
	if len(explicit) == 0 {
		// Explicit but empty — same fallthrough as null; treat as env-sourced.
		return envVal, false
	}
	return explicit, true
}

// dropCrossFlowEnvValues drops env-sourced auth values that belong to a
// different flow than the one the user picked via explicit attributes.
// Schema ConflictsWith already covers explicit-vs-explicit; this covers
// explicit-vs-env in either direction (e.g. explicit assertion + env
// BTP_USERNAME, or explicit username + env BTP_ASSERTION).
// x509 has no env fallback, so its explicit fields only participate as
// flow selectors.
// ponytail: idp is a modifier (works with any flow), not touched here.
func dropCrossFlowEnvValues(
	cfg providerData,
	username string, usernameExplicit bool,
	password string, passwordExplicit bool,
	idToken string, idTokenExplicit bool,
	assertion string, assertionExplicit bool,
	resp *provider.ConfigureResponse,
) (string, string, string, string) {
	explicitUserPw := usernameExplicit || passwordExplicit
	explicitX509 := !cfg.TLSClientKey.IsNull() || !cfg.TLSClientCertificate.IsNull() || !cfg.IdentityProviderURL.IsNull()

	warn := func(envName string) {
		resp.Diagnostics.AddWarning(
			"Conflicting authentication configuration",
			fmt.Sprintf("An authentication flow was selected via explicit provider attributes, but the environment variable %q is also set. The explicit configuration takes precedence; %q is ignored.", envName, envName),
		)
	}

	// user/pw flow not explicitly picked, but env values leak in while another
	// flow was picked explicitly — drop env-sourced user/pw.
	otherThanUserPw := idTokenExplicit || assertionExplicit || explicitX509
	if !explicitUserPw && otherThanUserPw {
		if !usernameExplicit && len(username) > 0 {
			warn("BTP_USERNAME")
			username = ""
		}
		if !passwordExplicit && len(password) > 0 {
			warn("BTP_PASSWORD")
			password = ""
		}
	}
	otherThanIdToken := explicitUserPw || assertionExplicit || explicitX509
	if !idTokenExplicit && otherThanIdToken && len(idToken) > 0 {
		warn("BTP_IDTOKEN")
		idToken = ""
	}
	otherThanAssertion := explicitUserPw || idTokenExplicit || explicitX509
	if !assertionExplicit && otherThanAssertion && len(assertion) > 0 {
		warn("BTP_ASSERTION")
		assertion = ""
	}
	return username, password, idToken, assertion
}

func determineAuthFlow(config providerData, idToken string, ssoLogin bool, assertion string, btpCliSessionLogin bool) string {
	if ssoLogin {
		return ssoFlow
	} else if len(idToken) > 0 {
		return idTokenFlow
	} else if !config.TLSClientKey.IsNull() {
		return x509Flow
	} else if len(assertion) > 0 {
		return assertionFlow
	} else if btpCliSessionLogin {
		return btpCliSessionFlow
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
