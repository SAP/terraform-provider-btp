package btpcli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestV2Client_New(t *testing.T) {
	t.Parallel()

	fakeURL, _ := url.Parse("https://my.cli.server.local")
	t.Run("correlationID generator correctly set", func(t *testing.T) {
		uut := NewV2Client(fakeURL)

		correlationID1 := uut.newCorrelationID()
		correlationID2 := uut.newCorrelationID()

		assert.NotEmpty(t, correlationID1)
		assert.NotEmpty(t, correlationID2)
		assert.NotEqual(t, correlationID1, correlationID2)
	})
	t.Run("default http client set", func(t *testing.T) {
		uut := NewV2Client(fakeURL)

		assert.NotNil(t, uut.httpClient)
	})
	t.Run("server url correctly set", func(t *testing.T) {
		uut := NewV2Client(fakeURL)

		assert.Equal(t, fakeURL, uut.serverURL)
	})
}

func TestV2Client_ProtocolVersion(t *testing.T) {
	assert.Regexp(t, regexp.MustCompile(`^v\d+\.\d+\.\d+$`), cliTargetProtocolVersion, "cliTargetProtocolVersion must be valid semver")
}

func TestV2Client_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string

		loginRequest *LoginRequest
		simulation   v2SimulationConfig
	}{
		{
			description:  "happy path",
			loginRequest: NewLoginRequest("subdomain", "john.doe", "pass"),
			simulation: v2SimulationConfig{
				srvExpectBody:    `{"customIdp":"","subdomain":"subdomain","userName":"john.doe","password":"pass"}`,
				srvReturnStatus:  http.StatusOK,
				srvReturnContent: `{"issuer": "accounts.sap.com","user":"john.doe","mail":"john.doe@test.com","refreshToken":"abc"}`,
				expectResponse: &LoginResponse{
					Issuer:       "accounts.sap.com",
					Username:     "john.doe",
					Email:        "john.doe@test.com",
					RefreshToken: "abc",
				},
				expectClientSession: &Session{
					RefreshToken:           "abc",
					IdentityProvider:       "",
					GlobalAccountSubdomain: "subdomain",
					LoggedInUser: &v2LoggedInUser{
						Issuer:   "accounts.sap.com",
						Username: "john.doe",
						Email:    "john.doe@test.com",
					},
				},
			},
		},
		{
			description:  "happy path - with custom idp",
			loginRequest: NewLoginRequestWithCustomIDP("my.custom.idp", "subdomain", "john.doe", "pass"),
			simulation: v2SimulationConfig{
				srvExpectBody:    `{"customIdp":"my.custom.idp","subdomain":"subdomain","userName":"john.doe","password":"pass"}`,
				srvReturnStatus:  http.StatusOK,
				srvReturnContent: `{"issuer": "customidp.accounts.ondemand.com","user":"john.doe","mail":"john.doe@test.com","refreshToken":"abc"}`,
				expectResponse: &LoginResponse{
					Issuer:       "customidp.accounts.ondemand.com",
					Username:     "john.doe",
					Email:        "john.doe@test.com",
					RefreshToken: "abc",
				},
				expectClientSession: &Session{
					RefreshToken:           "abc",
					GlobalAccountSubdomain: "subdomain",
					IdentityProvider:       "my.custom.idp",
					LoggedInUser: &v2LoggedInUser{
						Issuer:   "customidp.accounts.ondemand.com",
						Username: "john.doe",
						Email:    "john.doe@test.com",
					},
				},
			},
		},
		{
			description:  "error path - wrong credentials [401]",
			loginRequest: NewLoginRequest("subdomain", "john.doe", "this.is.wrong"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusUnauthorized,
				expectErrorMsg:  "Login failed. Check your credentials. [Status: 401; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - user is lacking permissions to globalaccount [403]",
			loginRequest: NewLoginRequest("subdomain", "john.doe", "pass"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusForbidden,
				expectErrorMsg:  "You cannot access global account 'subdomain'. Make sure you have at least read access to the global account, a directory, or a subaccount. [Status: 403; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - global account can't be found [404]",
			loginRequest: NewLoginRequest("subdomain", "john.doe", "pass"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusNotFound,
				expectErrorMsg:  "Global account 'subdomain' not found. Try again and make sure to provide the global account's subdomain. [Status: 404; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - login request times out [504]]",
			loginRequest: NewLoginRequest("subdomain", "john.doe", "pass"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusGatewayTimeout,
				expectErrorMsg:  "Login timed out. Please try again later. [Status: 504; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - unexpected error",
			loginRequest: NewLoginRequest("subdomain", "john.doe", "pass"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusTeapot,
				expectErrorMsg:  "Received response with unexpected status [Status: 418; Correlation ID: fake-correlation-id]",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.simulation.srvExpectPath = path.Join("/login", cliTargetProtocolVersion)
			test.simulation.callFunctionUnderTest = func(ctx context.Context, uut *v2Client) (any, error) {
				return uut.Login(ctx, test.loginRequest)
			}

			simulateV2Call(t, test.simulation)
		})
	}
}

func TestV2Client_Logout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string

		logoutRequest *LogoutRequest
		simulation    v2SimulationConfig
	}{
		{
			description:   "happy path",
			logoutRequest: NewLogoutRequest("subdomain"),
			simulation: v2SimulationConfig{
				srvExpectBody:    `{"customIdp":"","subdomain":"subdomain","refreshToken":""}`,
				srvReturnStatus:  http.StatusOK,
				srvReturnContent: `{"issuer": "","user":"john.doe","mail":"john.doe@test.com","refreshToken":"abc"}`,
				expectResponse:   &LogoutResponse{},
			},
		},
		{
			description:   "error path - login request times out [504]]",
			logoutRequest: NewLogoutRequest("subdomain"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusGatewayTimeout,
				expectErrorMsg:  "Logout timed out. Please try again later. [Status: 504; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:   "error path - unexpected error",
			logoutRequest: NewLogoutRequest("subdomain"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusTeapot,
				expectErrorMsg:  "Received response with unexpected status [Status: 418; Correlation ID: fake-correlation-id]",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.simulation.srvExpectPath = path.Join("/logout", cliTargetProtocolVersion)
			test.simulation.callFunctionUnderTest = func(ctx context.Context, uut *v2Client) (any, error) {
				return uut.Logout(ctx, test.logoutRequest)
			}

			simulateV2Call(t, test.simulation)
		})
	}
}

func TestV2Client_Execute(t *testing.T) {
	t.Run("every request must have a unique correlation ID", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !assert.Equal(t, path.Join("/command", cliTargetProtocolVersion, "subaccount/role"), r.URL.Path) {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set(HeaderCLIBackendStatus, "200")
			assertV2DefaultHeader(t, r, http.MethodPost)
			assert.Equal(t, string(ActionGet), r.URL.RawQuery)
			fmt.Fprintf(w, "{}")
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
		_, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.NoError(t, err)
	})
	t.Run("backend headers get passed through", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(HeaderCLIBackendStatus, fmt.Sprintf("%d", 201))
			w.Header().Set(HeaderCLIBackendMediaType, "backend/mediatype")
			fmt.Fprintf(w, "{}")
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
		res, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.NoError(t, err)
		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, "backend/mediatype", res.ContentType)
	})
	t.Run("custom idp: request header `X-CPCLI-CustomIdp` must be set", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "my.custom.idp", r.Header.Get(HeaderCLICustomIDP))
			w.Header().Set(HeaderCLIBackendStatus, fmt.Sprintf("%d", 201))
			w.Header().Set(HeaderCLIBackendMediaType, "backend/mediatype")
			fmt.Fprintf(w, "{}")
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
		uut.session = &Session{
			GlobalAccountSubdomain: "globalaccount-subdomain",
			IdentityProvider:       "my.custom.idp",
			LoggedInUser: &v2LoggedInUser{
				Email:    "john.doe@int.test",
				Username: "john.doe@int.test",
				Issuer:   "customidp.accounts.ondemand.com",
			},
		}

		_, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.NoError(t, err)
	})
}

type v2SimulationConfig struct {
	// initialize the client session prior to the test simulation
	initSession *Session

	// the api endpoint the client is expected to call
	srvExpectPath string

	// can be used to verify the payload sent by the client
	srvExpectBody string

	// the http status the fake api shall respond
	srvReturnStatus int

	// the content the fake api server shall respond
	srvReturnContent string

	// triggers the function under test on the uut
	callFunctionUnderTest func(ctx context.Context, uut *v2Client) (any, error)

	// expected error message (if any)
	expectErrorMsg string

	// expected client response
	expectResponse any

	// expected client session after the simulation
	expectClientSession *Session
}

func simulateV2Call(t *testing.T, config v2SimulationConfig) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !assert.Equal(t, config.srvExpectPath, r.URL.Path) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		b, err := io.ReadAll(r.Body)

		if assert.NoError(t, err) {
			assertV2DefaultHeader(t, r, http.MethodPost)

			if len(config.srvExpectBody) > 0 {
				assert.Equal(t, config.srvExpectBody, strings.TrimSpace(string(b)))
			}

			w.WriteHeader(config.srvReturnStatus)
			fmt.Fprintf(w, config.srvReturnContent)
		}
	}))
	defer srv.Close()

	srvUrl, _ := url.Parse(srv.URL)
	uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
	uut.session = config.initSession
	uut.newCorrelationID = func() string {
		return "fake-correlation-id"
	}
	response, err := config.callFunctionUnderTest(context.TODO(), uut)

	if len(config.expectErrorMsg) > 0 {
		assert.EqualError(t, err, config.expectErrorMsg)
	} else {
		if assert.NoError(t, err) && assert.NotNil(t, response) {
			assert.Equal(t, config.expectResponse, response)
		}
	}

	assert.Equal(t, config.expectClientSession, uut.session)
}

func assertV2DefaultHeader(t *testing.T, r *http.Request, expectedMethod string) {
	t.Helper()

	assert.Equal(t, expectedMethod, r.Method)
	assert.NotEmpty(t, r.Header.Get(HeaderCorrelationID))
	assert.Empty(t, r.Header.Get(HeaderCLICustomIDP))
	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, "json", r.Header.Get(HeaderCLIFormat))
}

func TestV2Client_GetLoggedInUser(t *testing.T) {
	t.Parallel()
	t.Run("no one logged in so far", func(t *testing.T) {
		uut := NewV2Client(nil)
		assert.Nil(t, uut.GetLoggedInUser())
	})
	t.Run("someone logged in", func(t *testing.T) {
		testUser := &v2LoggedInUser{
			Email:    "john.doe@test.com",
			Username: "john.doe",
			Issuer:   "test",
		}

		uut := NewV2Client(nil)
		uut.session = &Session{LoggedInUser: testUser}

		assert.Equal(t, testUser, uut.GetLoggedInUser())
	})
}

func TestV2Client_GetGlobalAccountSubdomain(t *testing.T) {
	t.Parallel()
	t.Run("no one logged in so far", func(t *testing.T) {
		uut := NewV2Client(nil)
		assert.Empty(t, uut.GetGlobalAccountSubdomain())
	})
	t.Run("someone logged in", func(t *testing.T) {
		uut := NewV2Client(nil)
		uut.session = &Session{GlobalAccountSubdomain: "my-subdomain"}

		assert.Equal(t, "my-subdomain", uut.GetGlobalAccountSubdomain())
	})
}
