package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	testingResource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"

	"github.com/SAP/terraform-provider-btp/internal/btpcli"
	"github.com/SAP/terraform-provider-btp/internal/btpcli/types/servicemanager"
)

// ─────────────────────── dummy data ────────────────────────────────────────

const (
	unitTestBrokerSubaccountID = "6aa64c2f-38c1-49a9-b2e8-cf9fea769b7f"
	unitTestBrokerID           = "9ff44f1b-b2a8-43ae-9072-32bd1dce60e4"
	unitTestBrokerName         = "unit-test-broker"
	unitTestBrokerURL          = "https://unit-test-broker.example.com"
	unitTestBrokerDesc         = "Unit test service broker"
	unitTestBrokerUser         = "broker-user"
	unitTestBrokerPwd          = "broker-password"
)

func newDummyBroker() servicemanager.ServiceBrokerResponseObject {
	return servicemanager.ServiceBrokerResponseObject{
		Id:          unitTestBrokerID,
		Name:        unitTestBrokerName,
		Description: unitTestBrokerDesc,
		BrokerUrl:   unitTestBrokerURL,
		Ready:       true,
		LastOperation: &servicemanager.OperationResponseObject{
			State: servicemanager.StateSucceeded,
			Type_: servicemanager.TypeCreate,
		},
		CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	}
}

// ─────────────────────── fake HTTP servers ──────────────────────────────────

// newBrokerFakeServer returns a TLS test server that handles login and all
// services/broker operations, returning a static successful response.
func newBrokerFakeServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/login/v2.106.1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"email":  "unit-test@example.com",
			"issuer": "accounts.sap.com",
		})
	})

	mux.HandleFunc("/command/v2.106.1/services/broker", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set(btpcli.HeaderCLIBackendStatus, "200")

		action := r.URL.RawQuery
		switch action {
		case "register", "get", "update":
			_ = json.NewEncoder(w).Encode(newDummyBroker())
		case "unregister":
			_ = json.NewEncoder(w).Encode(map[string]string{})
		default:
			http.Error(w, "unknown action", http.StatusBadRequest)
		}
	})

	return httptest.NewTLSServer(mux)
}

// newBrokerFakeServerWithRegisterError returns a TLS test server that succeeds
// on login but returns an HTTP 500 error for any broker API call.
func newBrokerFakeServerWithRegisterError(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/login/v2.106.1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"email":  "unit-test@example.com",
			"issuer": "accounts.sap.com",
		})
	})

	mux.HandleFunc("/command/v2.106.1/services/broker", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Signal backend error via header; HTTP 200 at transport level prevents go-retryablehttp retries.
		w.Header().Set(btpcli.HeaderCLIBackendStatus, "500")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "service broker registration failed",
		})
	})

	return httptest.NewTLSServer(mux)
}

// ─────────────────────── HCL helpers ────────────────────────────────────────

func hclBrokerResourceWithUsernamePassword(subaccountID, password string) string {
	return fmt.Sprintf(`
resource "btp_subaccount_service_broker" "test" {
  subaccount_id = %q
  name          = %q
  url           = %q
  username      = %q
  password      = %q
}
`, subaccountID, unitTestBrokerName, unitTestBrokerURL, unitTestBrokerUser, password)
}

// ─────────────────────── DIRECT UNIT TESTS (no HTTP) ────────────────────────

func TestSubaccountServiceBrokerMetadata(t *testing.T) {
	rs := newSubaccountServiceBrokerResource().(*subaccountServiceBrokerResource)

	req := fwresource.MetadataRequest{ProviderTypeName: "btp"}
	resp := &fwresource.MetadataResponse{}
	rs.Metadata(context.Background(), req, resp)

	assert.Equal(t, "btp_subaccount_service_broker", resp.TypeName)
}

func TestSubaccountServiceBrokerSchema(t *testing.T) {
	rs := newSubaccountServiceBrokerResource().(*subaccountServiceBrokerResource)

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	rs.Schema(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, int64(1), resp.Schema.Version)

	attrs := resp.Schema.Attributes

	// All 12 expected attributes
	for _, name := range []string{
		"subaccount_id", "id", "name", "description", "url",
		"username", "password", "cert", "key",
		"created_date", "last_modified", "ready",
	} {
		assert.Contains(t, attrs, name, "missing attribute: %s", name)
	}

	assert.True(t, attrs["subaccount_id"].IsRequired(), "subaccount_id should be required")
	assert.True(t, attrs["name"].IsRequired(), "name should be required")
	assert.True(t, attrs["url"].IsRequired(), "url should be required")

	assert.True(t, attrs["id"].IsComputed(), "id should be computed")
	assert.True(t, attrs["created_date"].IsComputed(), "created_date should be computed")
	assert.True(t, attrs["last_modified"].IsComputed(), "last_modified should be computed")
	assert.True(t, attrs["ready"].IsComputed(), "ready should be computed")

	assert.True(t, attrs["description"].IsOptional(), "description should be optional")
	assert.True(t, attrs["username"].IsOptional(), "username should be optional")
	assert.True(t, attrs["password"].IsOptional(), "password should be optional")
	assert.True(t, attrs["cert"].IsOptional(), "cert should be optional")
	assert.True(t, attrs["key"].IsOptional(), "key should be optional")

	// v1 schema must NOT contain the removed mtls attribute
	_, hasMtls := attrs["mtls"]
	assert.False(t, hasMtls, "v1 schema must not contain the mtls attribute")
}

func TestSubaccountServiceBrokerIdentitySchema(t *testing.T) {
	rs := newSubaccountServiceBrokerResource().(*subaccountServiceBrokerResource)

	req := fwresource.IdentitySchemaRequest{}
	resp := &fwresource.IdentitySchemaResponse{}
	rs.IdentitySchema(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, resp.IdentitySchema.Attributes, 2)
	assert.Contains(t, resp.IdentitySchema.Attributes, "subaccount_id")
	assert.Contains(t, resp.IdentitySchema.Attributes, "id")
}

func TestSubaccountServiceBrokerConfigure(t *testing.T) {
	t.Run("nil ProviderData leaves cli unset", func(t *testing.T) {
		rs := &subaccountServiceBrokerResource{}
		req := fwresource.ConfigureRequest{ProviderData: nil}
		resp := &fwresource.ConfigureResponse{}

		rs.Configure(context.Background(), req, resp)

		assert.Nil(t, rs.cli)
		assert.False(t, resp.Diagnostics.HasError())
	})
}

func TestSubaccountServiceBrokerValueFrom(t *testing.T) {
	createdAt := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 3, 16, 12, 0, 0, 0, time.UTC)

	broker := servicemanager.ServiceBrokerResponseObject{
		Id:          "broker-id-123",
		Name:        "my-broker",
		Description: "My test broker",
		BrokerUrl:   "https://broker.example.com",
		Ready:       true,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	result := subaccountServiceBrokerValueFrom(context.Background(), broker)

	assert.Equal(t, "broker-id-123", result.Id.ValueString())
	assert.Equal(t, "my-broker", result.Name.ValueString())
	assert.Equal(t, "My test broker", result.Description.ValueString())
	assert.Equal(t, "https://broker.example.com", result.Url.ValueString())
	assert.True(t, result.Ready.ValueBool())
	assert.False(t, result.CreatedDate.IsNull(), "CreatedDate should not be null")
	assert.False(t, result.LastModified.IsNull(), "LastModified should not be null")
	// SubaccountId, Username, Password, Cert, Key are intentionally zero—filled by caller
	assert.True(t, result.SubaccountId.IsNull())
	assert.True(t, result.Username.IsNull())
	assert.True(t, result.Password.IsNull())
	assert.True(t, result.Cert.IsNull())
	assert.True(t, result.Key.IsNull())
}

func TestSubaccountServiceBrokerImportState_WrongFormat(t *testing.T) {
	rs := &subaccountServiceBrokerResource{}

	cases := []struct {
		name string
		id   string
	}{
		{"single part", "just-an-id"},
		{"three parts", "a,b,c"},
		{"empty subaccount_id", ",broker-id"},
		{"empty broker id", "subaccount-id,"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := fwresource.ImportStateRequest{ID: tc.id}
			resp := &fwresource.ImportStateResponse{}

			rs.ImportState(context.Background(), req, resp)

			assert.True(t, resp.Diagnostics.HasError(), "expected error for import ID %q", tc.id)
		})
	}
}

func TestSubaccountServiceBrokerUpgradeState(t *testing.T) {
	rs := &subaccountServiceBrokerResource{}

	upgraders := rs.UpgradeState(context.Background())

	assert.Len(t, upgraders, 1, "expected exactly one state upgrader")

	upgrader, ok := upgraders[0]
	assert.True(t, ok, "state upgrader for version 0 must exist")
	assert.NotNil(t, upgrader.PriorSchema, "PriorSchema must not be nil")
	assert.NotNil(t, upgrader.StateUpgrader, "StateUpgrader func must not be nil")

	// v0 schema has the mtls bool attribute that was removed in v1
	_, hasMtls := upgrader.PriorSchema.Attributes["mtls"]
	assert.True(t, hasMtls, "v0 PriorSchema must contain the mtls attribute")

	// v0 schema must contain all base attributes that carry over to v1
	for _, name := range []string{"subaccount_id", "id", "name", "description", "url",
		"username", "password", "cert", "key", "created_date", "last_modified", "ready"} {
		assert.Contains(t, upgrader.PriorSchema.Attributes, name, "v0 schema missing attribute: %s", name)
	}
}

// ─────────────────────── TERRAFORM FRAMEWORK TESTS (fake server) ────────────

// TestSubaccountServiceBrokerCreate_Success tests a full Create → Read → Delete cycle.
func TestSubaccountServiceBrokerCreate_Success(t *testing.T) {
	srv := newBrokerFakeServer(t)
	defer srv.Close()

	testingResource.Test(t, testingResource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getProviders(srv.Client()),
		Steps: []testingResource.TestStep{
			{
				Config: hclProviderForCLIServerAt(srv.URL) +
					hclBrokerResourceWithUsernamePassword(unitTestBrokerSubaccountID, unitTestBrokerPwd),
				Check: testingResource.ComposeAggregateTestCheckFunc(
					testingResource.TestCheckResourceAttr("btp_subaccount_service_broker.test", "id", unitTestBrokerID),
					testingResource.TestCheckResourceAttr("btp_subaccount_service_broker.test", "name", unitTestBrokerName),
					testingResource.TestCheckResourceAttr("btp_subaccount_service_broker.test", "url", unitTestBrokerURL),
					testingResource.TestCheckResourceAttr("btp_subaccount_service_broker.test", "subaccount_id", unitTestBrokerSubaccountID),
					testingResource.TestCheckResourceAttr("btp_subaccount_service_broker.test", "ready", "true"),
					testingResource.TestCheckResourceAttr("btp_subaccount_service_broker.test", "description", unitTestBrokerDesc),
				),
			},
		},
	})
}

// TestSubaccountServiceBrokerUpdate_Success creates a broker and then updates
// its password, exercising the Update code path.
func TestSubaccountServiceBrokerUpdate_Success(t *testing.T) {
	srv := newBrokerFakeServer(t)
	defer srv.Close()

	testingResource.Test(t, testingResource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getProviders(srv.Client()),
		Steps: []testingResource.TestStep{
			{
				Config: hclProviderForCLIServerAt(srv.URL) +
					hclBrokerResourceWithUsernamePassword(unitTestBrokerSubaccountID, unitTestBrokerPwd),
				Check: testingResource.TestCheckResourceAttr(
					"btp_subaccount_service_broker.test", "id", unitTestBrokerID,
				),
			},
			{
				// Change password to trigger an Update
				Config: hclProviderForCLIServerAt(srv.URL) +
					hclBrokerResourceWithUsernamePassword(unitTestBrokerSubaccountID, "updated-password"),
				Check: testingResource.TestCheckResourceAttr(
					"btp_subaccount_service_broker.test", "id", unitTestBrokerID,
				),
			},
		},
	})
}

// TestSubaccountServiceBrokerImportState_Success creates a broker and then
// imports it using the subaccount_id,id format.
func TestSubaccountServiceBrokerImportState_Success(t *testing.T) {
	srv := newBrokerFakeServer(t)
	defer srv.Close()

	testingResource.Test(t, testingResource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getProviders(srv.Client()),
		Steps: []testingResource.TestStep{
			{
				Config: hclProviderForCLIServerAt(srv.URL) +
					hclBrokerResourceWithUsernamePassword(unitTestBrokerSubaccountID, unitTestBrokerPwd),
			},
			{
				ResourceName:      "btp_subaccount_service_broker.test",
				ImportState:       true,
				ImportStateVerify: false, // fake server returns static data; sensitive fields absent
				ImportStateId:     fmt.Sprintf("%s,%s", unitTestBrokerSubaccountID, unitTestBrokerID),
			},
		},
	})
}

// TestSubaccountServiceBrokerImportState_ErrorWrongFormat creates a broker and
// then attempts an import with a malformed ID, expecting a diagnostic error.
func TestSubaccountServiceBrokerImportState_ErrorWrongFormat(t *testing.T) {
	srv := newBrokerFakeServer(t)
	defer srv.Close()

	testingResource.Test(t, testingResource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getProviders(srv.Client()),
		Steps: []testingResource.TestStep{
			{
				Config: hclProviderForCLIServerAt(srv.URL) +
					hclBrokerResourceWithUsernamePassword(unitTestBrokerSubaccountID, unitTestBrokerPwd),
			},
			{
				ResourceName:  "btp_subaccount_service_broker.test",
				ImportState:   true,
				ImportStateId: "wrong-format-no-comma",
				ExpectError:   regexp.MustCompile(`Unexpected Import Identifier`),
			},
		},
	})
}

// TestSubaccountServiceBrokerCreate_ApiError verifies that when the BTP API
// returns an error during Register, the provider surfaces a diagnostic.
func TestSubaccountServiceBrokerCreate_ApiError(t *testing.T) {
	srv := newBrokerFakeServerWithRegisterError(t)
	defer srv.Close()

	testingResource.Test(t, testingResource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: getProviders(srv.Client()),
		Steps: []testingResource.TestStep{
			{
				Config: hclProviderForCLIServerAt(srv.URL) +
					hclBrokerResourceWithUsernamePassword(unitTestBrokerSubaccountID, unitTestBrokerPwd),
				ExpectError: regexp.MustCompile(`API Error Creating Resource Service Broker`),
			},
		},
	})
}
