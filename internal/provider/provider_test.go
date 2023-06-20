package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/SAP/terraform-provider-btp/internal/validation/uuidvalidator"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

var (
	regexpValidRFC3999Format = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	regexpValidUUID          = uuidvalidator.UuidRegexp
)

func hclProvider() string {
	return hclProviderWithCLIServerURL("https://cpcli.cf.sap.hana.ondemand.com")
}

func hclProviderWithCLIServerURL(cliServerURL string) string {
	// TODO replace credentials with serviceuser credentials
	return fmt.Sprintf(`
provider "btp" {
    cli_server_url = "%s"
    globalaccount  = "terraformintcanary"
    username       = "john.doe@int.test"
    password       = "redacted"
}
    `, cliServerURL)
}

func getProviders(httpClient *http.Client) map[string]func() (tfprotov6.ProviderServer, error) {
	btpProvider := NewWithClient(httpClient).(*btpcliProvider)
	btpProvider.betaFeaturesEnabled = true // allows beta resources/datasource to be int. tested
	return map[string]func() (tfprotov6.ProviderServer, error){
		"btp": providerserver.NewProtocol6WithError(btpProvider),
	}
}

func setupVCR(t *testing.T, cassetteName string) *recorder.Recorder {
	t.Helper()

	rec, err := recorder.New(cassetteName)

	if err != nil {
		t.Fatal()
	}
	hookRedactIntegrationUserCredentials := func(i *cassette.Interaction) error {
		intUser := os.Getenv("BTP_USERNAME")
		intUserPwd := os.Getenv("BTP_PASSWORD")

		firstName, lastName := getNameFromEmail(intUser)

		if strings.Contains(i.Request.URL, "/login/") {
			i.Request.Body = strings.ReplaceAll(i.Request.Body, intUserPwd, "redacted")
		}

		i.Request.Body = strings.ReplaceAll(i.Request.Body, intUser, "john.doe@int.test")
		i.Response.Body = strings.ReplaceAll(i.Response.Body, intUser, "john.doe@int.test")

		if strings.Contains(i.Response.Body, "givenName") {
			i.Response.Body = strings.ReplaceAll(i.Response.Body, firstName, "John")
		}

		if strings.Contains(i.Response.Body, "familyName") {
			i.Response.Body = strings.ReplaceAll(i.Response.Body, lastName, "Doe")
		}

		if strings.Contains(i.Response.Body, "externalId") {
			indexOfExternalId := strings.Index(i.Response.Body, "\"externalId\":")
			i.Response.Body = i.Response.Body[:indexOfExternalId+14] + "I000000" + i.Response.Body[indexOfExternalId+21:]
		}

		return nil
	}

	hookRedactTokensInHeader := func(i *cassette.Interaction) error {
		redactTokenHeaders := func(headers map[string][]string) {
			for key := range headers {
				if strings.Contains(strings.ToLower(key), "token") {
					headers[key] = []string{"redacted"}
				}
			}
		}

		redactTokenHeaders(i.Request.Headers)
		redactTokenHeaders(i.Response.Headers)

		re := regexp.MustCompile(`"refreshToken":\s*"([a-f0-9]+)"`)
		i.Request.Body = re.ReplaceAllString(i.Request.Body, `"refreshToken":"redacted"`)
		i.Response.Body = re.ReplaceAllString(i.Response.Body, `"refreshToken":"redacted"`)

		return nil
	}

	rec.AddHook(hookRedactIntegrationUserCredentials, recorder.BeforeSaveHook)
	rec.AddHook(hookRedactTokensInHeader, recorder.BeforeSaveHook)

	return rec
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

func TestProvider_HasResources(t *testing.T) {
	expectedResources := []string{
		"btp_directory",
		/* TODO: switched off for phase 1
		/"btp_directory_role",
		*/
		"btp_directory_role_collection",
		"btp_directory_role_collection_assignment",
		"btp_globalaccount_resource_provider",
		/* TODO: switched off for phase 1
		"btp_globalaccount_role",
		*/
		"btp_globalaccount_role_collection",
		"btp_globalaccount_role_collection_assignment",
		"btp_globalaccount_trust_configuration",
		"btp_subaccount",
		"btp_subaccount_entitlement",
		"btp_subaccount_environment_instance",
		"btp_subaccount_role",
		"btp_subaccount_role_collection",
		"btp_subaccount_role_collection_assignment",
		/* TODO: switched off for phase 1
		"btp_subaccount_service_binding",
		"btp_subaccount_service_instance",
		"btp_subaccount_subscription",
		*/
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
		/*TODO: Switched off for phase 1
		"btp_directory_app",
		"btp_directory_apps",
		*/
		"btp_directory_entitlements",
		"btp_directory_labels",
		"btp_directory_role",
		"btp_directory_role_collection",
		"btp_directory_role_collections",
		"btp_directory_roles",
		"btp_directory_user",
		"btp_directory_users",
		"btp_globalaccount",
		/*TODO: Switched off for phase 1
		"btp_globalaccount_app",
		"btp_globalaccount_apps",
		*/
		"btp_globalaccount_entitlements",
		/*TODO: Switched off for phase 1
		"btp_globalaccount_resource_provider",
		"btp_globalaccount_resource_providers",
		*/
		"btp_globalaccount_role",
		"btp_globalaccount_role_collection",
		"btp_globalaccount_role_collections",
		"btp_globalaccount_roles",
		"btp_globalaccount_trust_configuration",
		"btp_globalaccount_trust_configurations",
		"btp_globalaccount_user",
		"btp_globalaccount_users",
		"btp_regions",
		"btp_subaccount",
		/*TODO: Switched off for phase 1
		"btp_subaccount_app",
		"btp_subaccount_apps",
		*/
		"btp_subaccount_entitlements",
		"btp_subaccount_environment_instance",
		"btp_subaccount_environment_instances",
		"btp_subaccount_environments",
		"btp_subaccount_labels",
		"btp_subaccount_role",
		"btp_subaccount_role_collection",
		"btp_subaccount_role_collections",
		"btp_subaccount_roles",
		/*TODO: Switched off for phase 1
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
		"btp_subaccount_service_platform",
		"btp_subaccount_service_platforms",
		"btp_subaccount_subscription",
		"btp_subaccount_subscriptions",
		*/
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
