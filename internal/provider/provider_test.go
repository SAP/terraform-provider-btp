package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	testingResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	regexpValidRFC3999Format = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	regexpValidUUID          = uuidvalidator.UuidRegexp
)

type TestUser struct {
	Username  string
	Password  string
	Idp       string
	Issuer    string
	Firstname string
	Lastname  string
}

var redactedTestUser = TestUser{
	Username:  "john.doe@int.test",
	Password:  "testUserPassword",
	Idp:       "identityProvider",
	Issuer:    "identity.provider.test",
	Firstname: "John",
	Lastname:  "Doe",
}

var testGlobalAccount = getenv("BTP_GLOBALACCOUNT", "terraformintcanary")

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func hclProviderFor(user TestUser) string {
	return hclProvider("https://canary.cli.btp.int.sap", user)
}

func hclProviderForCLIServerAt(cliServerURL string) string {
	return hclProvider(cliServerURL, redactedTestUser)
}

func hclProvider(cliServerURL string, user TestUser) string {
	// TODO replace credentials with serviceuser credentials
	return fmt.Sprintf(`
provider "btp" {
    cli_server_url = "%s"
    globalaccount  = "%s"
    username       = "%s"
    password       = "%s"
    idp            = "%s"
}
    `, cliServerURL, testGlobalAccount, user.Username, user.Password, user.Idp)
}

func getProviders(httpClient *http.Client) map[string]func() (tfprotov6.ProviderServer, error) {
	btpProvider := NewWithClient(httpClient).(*btpcliProvider)
	btpProvider.betaFeaturesEnabled = true // allows beta resources/datasource to be int. tested
	return map[string]func() (tfprotov6.ProviderServer, error){
		"btp": providerserver.NewProtocol6WithError(btpProvider),
	}
}

func setupVCR(t *testing.T, cassetteName string) (*recorder.Recorder, TestUser) {
	t.Helper()

	mode := recorder.ModeRecordOnce
	if force, _ := strconv.ParseBool(os.Getenv("TEST_FORCE_REC")); force {
		mode = recorder.ModeRecordOnly
	}

	rec, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       cassetteName,
		Mode:               mode,
		SkipRequestLatency: true,
		RealTransport:      http.DefaultTransport,
	})

	user := redactedTestUser
	if rec.IsRecording() {
		t.Logf("ATTENTION: Recording '%s'", cassetteName)
		user.Username = os.Getenv("BTP_USERNAME")
		user.Password = os.Getenv("BTP_PASSWORD")
		if len(user.Username) == 0 || len(user.Password) == 0 {
			t.Fatal("Env vars BTP_USERNAME and BTP_PASSWORD are required when recording test fixtures")
		}

		user.Idp = os.Getenv("BTP_IDP")
		if len(user.Idp) == 0 {
			user.Issuer = "accounts.sap.com"
		} else if strings.Contains(user.Idp, ".") {
			user.Issuer = user.Idp
		} else {
			// TODO: currently short notation (idp = tenant id) only supported with this default domain (e.g. canary)
			user.Issuer = user.Idp + ".accounts400.ondemand.com"
		}
		user.Firstname, user.Lastname = getNameFromEmail(user.Username)
	} else {
		t.Logf("Replaying '%s'", cassetteName)
	}

	if err != nil {
		t.Fatal()
	}

	rec.SetMatcher(cliServerRequestMatcher(t))
	rec.AddHook(hookRedactIntegrationUserCredentials(user), recorder.BeforeSaveHook)
	rec.AddHook(hookRedactTokensAndSessionIds(), recorder.BeforeSaveHook)

	return rec, user
}

func cliServerRequestMatcher(t *testing.T) func(r *http.Request, i cassette.Request) bool {
	return func(r *http.Request, i cassette.Request) bool {
		if r.Method != i.Method || r.URL.String() != i.URL {
			return false
		}

		subdomainHeaderKey := http.CanonicalHeaderKey(btpcli.HeaderCLISubdomain)
		if r.Header.Get(subdomainHeaderKey) != i.Headers.Get(subdomainHeaderKey) {
			return false
		}

		idpHeaderKey := http.CanonicalHeaderKey(btpcli.HeaderCLICustomIDP)
		if r.Header.Get(idpHeaderKey) != i.Headers.Get(idpHeaderKey) {
			return false
		}

		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal("Unable to read body from request")
		}
		requestBody := string(bytes)
		return requestBody == i.Body
	}
}

func hookRedactIntegrationUserCredentials(user TestUser) func(i *cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
		if strings.Contains(i.Request.URL, "/login/") {
			reUserPwd := regexp.MustCompile(`"password":"(.*?)"`)
			i.Request.Body = reUserPwd.ReplaceAllString(i.Request.Body, `"password":"`+redactedTestUser.Password+`"`)
			reCustomIdp := regexp.MustCompile(`"customIdp":"(.*?)"`)
			i.Request.Body = reCustomIdp.ReplaceAllString(i.Request.Body, `"customIdp":"`+redactedTestUser.Idp+`"`)
			reIssuer := regexp.MustCompile(`"issuer":"(.*?)"`)
			i.Response.Body = reIssuer.ReplaceAllString(i.Response.Body, `"issuer":"`+redactedTestUser.Issuer+`"`)
		}

		if _, exists := i.Request.Headers[btpcli.HeaderCLICustomIDP]; exists {
			i.Request.Headers.Set(btpcli.HeaderCLICustomIDP, redactedTestUser.Idp)
		}

		reUserSAP := regexp.MustCompile(`[a-zA-Z0-9.!#$%&'*+/=?^_ \x60{|}~-]+@sap\.com`)
		i.Request.Body = reUserSAP.ReplaceAllString(i.Request.Body, redactedTestUser.Username)
		i.Response.Body = strings.ReplaceAll(i.Response.Body, user.Username, redactedTestUser.Username)
		// to support responses containing sets of email addresses, we need to replace with unique values
		index := 0
		i.Response.Body = reUserSAP.ReplaceAllStringFunc(i.Response.Body, func(string) string {
			index++
			return strings.ReplaceAll(redactedTestUser.Username, "@", "+"+strconv.Itoa(index)+"@")
		})

		if strings.Contains(i.Response.Body, "givenName") {
			i.Response.Body = strings.ReplaceAll(i.Response.Body, user.Firstname, redactedTestUser.Firstname)
		}

		if strings.Contains(i.Response.Body, "familyName") {
			i.Response.Body = strings.ReplaceAll(i.Response.Body, user.Lastname, redactedTestUser.Lastname)
		}

		if strings.Contains(i.Response.Body, "externalId") {
			indexOfExternalId := strings.Index(i.Response.Body, "\"externalId\":")
			i.Response.Body = i.Response.Body[:indexOfExternalId+14] + "I000000" + i.Response.Body[indexOfExternalId+21:]
		}

		if strings.Contains(i.Response.Body, "clientid") {
			reClientSecretVariant1 := regexp.MustCompile(`"clientid":"(.*?)"`)
			i.Response.Body = reClientSecretVariant1.ReplaceAllString(i.Response.Body, `"clientid":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "clientsecret") {
			reClientSecretVariant1 := regexp.MustCompile(`"clientsecret":"(.*?)"`)
			i.Response.Body = reClientSecretVariant1.ReplaceAllString(i.Response.Body, `"clientsecret":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "client_id") {
			reClientSecretVariant2 := regexp.MustCompile(`"client_id":"(.*?)"`)
			i.Response.Body = reClientSecretVariant2.ReplaceAllString(i.Response.Body, `"client_id":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "client_secret") {
			reClientSecretVariant2 := regexp.MustCompile(`"client_secret":"(.*?)"`)
			i.Response.Body = reClientSecretVariant2.ReplaceAllString(i.Response.Body, `"client_secret":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "username") {
			reBindingSecret := regexp.MustCompile(`"username":"(.*?)"`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `"username":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "password") {
			reBindingSecret := regexp.MustCompile(`"password":"(.*?)"`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `"password":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "key") {
			reBindingSecret := regexp.MustCompile(`"key":"(.*?)"`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `"key":"redacted"`)
		}

		if strings.Contains(i.Response.Body, "certificate") {
			reBindingSecret := regexp.MustCompile(`"certificate":"(.*?)"`)
			i.Response.Body = reBindingSecret.ReplaceAllString(i.Response.Body, `"certificate":"redacted"`)
		}

		if strings.Contains(i.Request.Body, "certificate") {
			reBindingSecret := regexp.MustCompile(`"certificate":"(.*?)"`)
			i.Request.Body = reBindingSecret.ReplaceAllString(i.Request.Body, `"certificate":"redacted"`)
		}

		return nil
	}
}

func hookRedactTokensAndSessionIds() func(i *cassette.Interaction) error {
	return func(i *cassette.Interaction) error {
		redact := func(headers map[string][]string) {
			for key := range headers {
				if strings.Contains(strings.ToLower(key), "token") ||
					strings.Contains(strings.ToLower(key), "session") {
					headers[key] = []string{"redacted"}
				}
			}
		}

		redact(i.Request.Headers)
		redact(i.Response.Headers)

		// TODO: this can be removed as soon as the btp CLI server no longer
		//  includes the refresh token in the body for compatibility reasons
		re := regexp.MustCompile(`"refreshToken":"(.*?)"`)
		i.Request.Body = re.ReplaceAllString(i.Request.Body, `"refreshToken":"redacted"`)
		i.Response.Body = re.ReplaceAllString(i.Response.Body, `"refreshToken":"redacted"`)

		return nil
	}
}

func stopQuietly(rec *recorder.Recorder) {
	if err := rec.Stop(); err != nil {
		panic(err)
	}
}

func getNameFromEmail(email string) (firstName, lastName string) {
	emailAt := strings.Index(email, "@")
	emailFirstName := strings.Split(email[:emailAt], ".")[0]
	emailLastName := strings.Split(email[:emailAt], ".")[1]

	firstName = convertFirstLetterToUpperCase(emailFirstName)
	lastName = convertFirstLetterToUpperCase(emailLastName)
	return
}

func convertFirstLetterToUpperCase(stringToConvert string) (convertedString string) {
	runes := []rune(stringToConvert)
	runes[0] = unicode.ToUpper(runes[0])
	convertedString = string(runes)
	return
}

func containsCheckFunc(expectedSubString string) testingResource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if !strings.Contains(value, expectedSubString) {
			return fmt.Errorf("expected value containing '%s', got: %s", expectedSubString, value)
		}
		return nil
	}
}

func notContainsCheckFunc(unexpectedSubString string) testingResource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if strings.Contains(value, unexpectedSubString) {
			return fmt.Errorf("expected value NOT containing '%s', got: %s", unexpectedSubString, value)
		}
		return nil
	}
}

func TestProvider_ConfigurationFlows(t *testing.T) {
	t.Parallel()
	t.Run("error path - user password login with missing data", func(t *testing.T) {
		rec, _ := setupVCR(t, "fixtures/provider.error_user_pwd")
		defer stopQuietly(rec)

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []testingResource.TestStep{
				{
					Config: `
provider "btp" {
	globalaccount  = "ga"
	username       = ""
	password       = "password"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`empty value for the username`),
				},
				{
					Config: `
provider "btp" {
	globalaccount  = "ga"
	username       = "username"
	password       = ""
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`empty value for the password`),
				},
			},
		})
	})
	t.Run("error path - x509 with missing data", func(t *testing.T) {
		rec, _ := setupVCR(t, "fixtures/provider.error_x509")
		defer stopQuietly(rec)

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []testingResource.TestStep{
				{
					Config: `
provider "btp" {
	globalaccount          = "ga"
	username               = ""
	tls_client_key         = "tlsClientKey"
	tls_client_certificate = "tlsClientCertificate"
	tls_idp_url            = "idpUrl"
	idp                    = "idp"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`empty value for the username`),
				},
				{
					Config: `
provider "btp" {
	globalaccount          = "ga"
	username               = "username"
	tls_client_key         = ""
	tls_client_certificate = "tlsClientCertificate"
	tls_idp_url            = "idpUrl"
	idp                    = "idp"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`empty value for the tls_client_key`),
				},
				{
					Config: `
provider "btp" {
	globalaccount          = "ga"
	username               = "username"
	tls_client_key         = "tlsClientKey"
	tls_client_certificate = ""
	tls_idp_url            = "idpUrl"
	idp                    = "idp"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`empty value for the tls_client_certificate`),
				},
				{
					Config: `
provider "btp" {
	globalaccount          = "ga"
	username               = "username"
	tls_client_key         = "tlsClientKey"
	tls_client_certificate = "tlsClientCertificate"
	tls_idp_url            = ""
	idp                    = "idp"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`empty value for the tls_idp_url`),
				},
				{
					Config: `
provider "btp" {
	globalaccount          = "ga"
	username               = "username"
	tls_client_key         = "tlsClientKey"
	tls_client_certificate = "tlsClientCertificate"
	tls_idp_url            = "idpUrl"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`Attribute "idp" must be specified`),
				},
			},
		})
	})
	t.Run("error path - invalid client server url", func(t *testing.T) {
		rec, _ := setupVCR(t, "fixtures/provider.error_cli_server_url")
		defer stopQuietly(rec)

		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(rec.GetDefaultClient()),
			Steps: []testingResource.TestStep{
				{
					Config: `
					provider "btp" {
						cli_server_url = "://canary.cli .btp.int.sap"
						globalaccount  = "ga"
						username       = "username"
						password       = "password"
					}
					data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`unableToCreateClient`),
				},
			},
		})
	})
}

func TestProvider_ConfigurationWithIdToken(t *testing.T) {
	t.Run("error path - attribute conflicts with idtoken", func(t *testing.T) {
		testingResource.Test(t, testingResource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: getProviders(nil),
			Steps: []testingResource.TestStep{
				{
					Config: `
provider "btp" {
	globalaccount  = "ga"
	username       = "username"
	idtoken        = "idtoken"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`Attribute "username" cannot be specified when "idtoken" is specified`),
				},
				{
					Config: `
provider "btp" {
	globalaccount  = "ga"
	password       = "password"
	idtoken        = "idtoken"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`Attribute "password" cannot be specified when "idtoken" is specified`),
				},
				{
					Config: `
provider "btp" {
	globalaccount  = "ga"
	idp            = "idp"
	idtoken        = "idtoken"
}
data "btp_whoami" "me" {}`,
					ExpectError: regexp.MustCompile(`Attribute "idp" cannot be specified when "idtoken" is specified`),
				},
			},
		})
	})
}

func TestProvider_HasResources(t *testing.T) {
	expectedResources := []string{
		"btp_directory_api_credential",
		"btp_directory",
		"btp_directory_entitlement",
		"btp_directory_role",
		"btp_directory_role_collection",
		"btp_directory_role_collection_assignment",
		"btp_globalaccount_api_credential",
		"btp_globalaccount_resource_provider",
		"btp_globalaccount_role",
		"btp_globalaccount_role_collection",
		"btp_globalaccount_role_collection_assignment",
		"btp_globalaccount_security_settings",
		"btp_globalaccount_trust_configuration",
		"btp_subaccount_api_credential",
		"btp_subaccount",
		"btp_subaccount_entitlement",
		"btp_subaccount_environment_instance",
		"btp_subaccount_role",
		"btp_subaccount_role_collection",
		"btp_subaccount_role_collection_assignment",
		"btp_subaccount_security_settings",
		"btp_subaccount_service_instance",
		"btp_subaccount_service_binding",
		"btp_subaccount_service_broker",
		"btp_subaccount_subscription",
		"btp_subaccount_trust_configuration",
	}

	ctx := context.Background()
	registeredResources := []string{}

	for _, resourceFunc := range New().Resources(ctx) {
		var resp resource.MetadataResponse

		resourceFunc().Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "btp"}, &resp)

		registeredResources = append(registeredResources, resp.TypeName)
	}

	assert.ElementsMatch(t, expectedResources, registeredResources)
}

func TestProvider_HasDatasources(t *testing.T) {
	expectedDataSources := []string{
		"btp_directory",
		"btp_directories",
		/*
			"btp_directory_app",
			"btp_directory_apps",
		*/
		"btp_directory_entitlements",
		"btp_directory_entitlement",
		"btp_directory_labels",
		"btp_directory_role",
		"btp_directory_role_collection",
		"btp_directory_role_collections",
		"btp_directory_roles",
		"btp_directory_user",
		"btp_directory_users",
		"btp_globalaccount",
		/*
			"btp_globalaccount_app",
			"btp_globalaccount_apps",
		*/
		"btp_globalaccount_entitlements",
		/*
			"btp_globalaccount_resource_provider",
			"btp_globalaccount_resource_providers",
		*/
		"btp_globalaccount_role",
		"btp_globalaccount_role_collection",
		"btp_globalaccount_role_collections",
		"btp_globalaccount_roles",
		"btp_globalaccount_security_settings",
		"btp_globalaccount_trust_configuration",
		"btp_globalaccount_trust_configurations",
		"btp_globalaccount_user",
		"btp_globalaccount_users",
		"btp_globalaccount_with_hierarchy",
		"btp_regions",
		"btp_subaccount",
		"btp_subaccount_app",
		"btp_subaccount_apps",
		"btp_subaccount_entitlements",
		"btp_subaccount_entitlement_unique_identifier",
		"btp_subaccount_environment_instance",
		"btp_subaccount_environment_instances",
		"btp_subaccount_environments",
		"btp_subaccount_labels",
		"btp_subaccount_role",
		"btp_subaccount_role_collection",
		"btp_subaccount_role_collections",
		"btp_subaccount_roles",
		"btp_subaccount_security_settings",
		"btp_subaccount_service_binding",
		"btp_subaccount_service_bindings",
		"btp_subaccount_service_broker",
		"btp_subaccount_service_brokers",
		"btp_subaccount_service_instance",
		"btp_subaccount_service_instances",
		"btp_subaccount_service_offering",
		"btp_subaccount_service_offerings",
		"btp_subaccount_service_plan",
		"btp_subaccount_service_plans",
		/*
			"btp_subaccount_service_platform",
			"btp_subaccount_service_platforms",
		*/
		"btp_subaccount_subscription",
		"btp_subaccount_subscriptions",
		"btp_subaccount_trust_configuration",
		"btp_subaccount_trust_configurations",
		"btp_subaccount_user",
		"btp_subaccount_users",
		"btp_subaccounts",
		"btp_whoami",
	}

	ctx := context.Background()
	registeredDataSources := []string{}

	for _, resourceFunc := range New().DataSources(ctx) {
		var resp datasource.MetadataResponse

		resourceFunc().Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "btp"}, &resp)

		registeredDataSources = append(registeredDataSources, resp.TypeName)
	}

	assert.ElementsMatch(t, expectedDataSources, registeredDataSources)
}
