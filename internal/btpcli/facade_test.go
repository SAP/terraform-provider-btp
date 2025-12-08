package btpcli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareClientFacadeForTest(handleFn http.HandlerFunc) (*ClientFacade, *httptest.Server) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderCLIBackendStatus, "200")

		handleFn.ServeHTTP(w, r)
	}))
	srvUrl, _ := url.Parse(srv.URL)

	apiClient := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
	apiClient.session = &Session{GlobalAccountSubdomain: "795b53bb-a3f0-4769-adf0-26173282a975"}
	return NewClientFacade(apiClient), srv
}

func assertCall(t *testing.T, r *http.Request, expectedCommand string, expectedAction Action, expectedParams map[string]string) {
	t.Helper()

	var payload struct {
		ParamValues map[string]string `json:"paramValues"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); assert.NoError(t, err) {
		expectedEndpoint := fmt.Sprintf("/command/%s/%s", cliTargetProtocolVersion, expectedCommand)

		assert.Equal(t, expectedEndpoint, r.URL.Path)
		assert.Equal(t, string(expectedAction), r.URL.RawQuery)
		assert.Equal(t, expectedParams, payload.ParamValues)
	}
}

func assertCallAnyMap(t *testing.T, r *http.Request, expectedCommand string, expectedAction Action, expectedParams map[string]any) {
	t.Helper()

	var payload struct {
		ParamValues map[string]any `json:"paramValues"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); assert.NoError(t, err) {
		expectedEndpoint := fmt.Sprintf("/command/%s/%s", cliTargetProtocolVersion, expectedCommand)

		assert.Equal(t, expectedEndpoint, r.URL.Path)
		assert.Equal(t, string(expectedAction), r.URL.RawQuery)
		assert.Equal(t, expectedParams, payload.ParamValues)
	}
}
