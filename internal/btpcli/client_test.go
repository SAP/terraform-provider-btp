package btpcli

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflogtest"
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

func TestV2Client_ServerVerboseEnvironment(t *testing.T) {
	t.Setenv(envBTPServerVerbose, "true")

	fakeURL, _ := url.Parse("https://my.cli.server.local")
	uut := NewV2Client(fakeURL)

	assert.True(t, uut.serverVerbose)
}

func TestV2Client_LogServerResponse(t *testing.T) {
	ctx := context.WithValue(context.Background(), v2ContextKey(HeaderCorrelationID), "correlation-id")
	errorBody := `{"error":"complete error response","details":["first","second"]}`
	response := &http.Response{
		StatusCode: http.StatusForbidden,
		Header: http.Header{
			HeaderCLIBackendStatus: []string{"401"},
			HeaderCLIServerMessage: []string{"session authorization failed"},
			HeaderCLISessionId:     []string{"must-not-be-logged"},
		},
		Body: io.NopCloser(strings.NewReader(errorBody)),
	}

	t.Run("disabled", func(t *testing.T) {
		called := false
		uut := NewV2Client(nil)
		uut.serverVerbose = false
		uut.verboseLogger = func(context.Context, serverVerboseEvent) {
			called = true
		}

		uut.logServerResponse(ctx, response)

		assert.False(t, called)
	})

	t.Run("enabled", func(t *testing.T) {
		var event serverVerboseEvent
		uut := NewV2Client(nil)
		uut.serverVerbose = true
		uut.verboseLogger = func(_ context.Context, received serverVerboseEvent) {
			event = received
		}
		uut.sessionEstablishedAt.Store(time.Now().Add(-2 * time.Minute).UnixNano())

		uut.logServerResponse(ctx, response)

		assert.Equal(t, "correlation-id", event.CorrelationID)
		assert.Equal(t, http.StatusForbidden, event.OuterHTTPStatus)
		assert.Equal(t, "401", event.BackendStatus)
		assert.Equal(t, "session authorization failed", event.ServerMessage)
		responseBody, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.Equal(t, errorBody, string(responseBody), "logging must not consume the response body")
		elapsed, err := time.ParseDuration(event.ElapsedSinceLogin)
		assert.NoError(t, err)
		assert.InDelta(t, (2 * time.Minute).Seconds(), elapsed.Seconds(), 0.1)
	})

	t.Run("backend error in successful outer response", func(t *testing.T) {
		backendErrorBody := `{"error":"backend error"}`
		response := &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				HeaderCLIBackendStatus: []string{"500"},
			},
			Body: io.NopCloser(strings.NewReader(backendErrorBody)),
		}

		uut := NewV2Client(nil)
		uut.serverVerbose = true
		uut.verboseLogger = func(context.Context, serverVerboseEvent) {}

		uut.logServerResponse(ctx, response)

		responseBody, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.Equal(t, backendErrorBody, string(responseBody), "logging must not consume backend error bodies")
	})

	t.Run("successful response body is not logged", func(t *testing.T) {
		response := &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				HeaderCLIBackendStatus: []string{"200"},
			},
			Body: io.NopCloser(strings.NewReader(`{"secret":"value"}`)),
		}

		uut := NewV2Client(nil)
		uut.serverVerbose = true
		uut.verboseLogger = func(context.Context, serverVerboseEvent) {}

		uut.logServerResponse(ctx, response)

		responseBody, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		assert.Equal(t, "{\"secret\":\"value\"}", string(responseBody), "logging must not consume successful response bodies")
	})
}

func TestV2Client_SessionRequestPolicies(t *testing.T) {
	tests := []struct {
		name                   string
		options                []V2ClientOption
		wantConcurrentRequests bool
	}{
		{
			name: "default serializes session requests",
		},
		{
			name:                   "concurrent option permits concurrent session requests",
			options:                []V2ClientOption{WithConcurrentSessionRequests()},
			wantConcurrentRequests: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstRequestMayComplete := make(chan struct{})
			requestReachedServer := make(chan struct{}, 2)
			requestHeaders := make(chan http.Header, 2)
			var requestCount atomic.Int32

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestHeaders <- r.Header.Clone()
				requestNumber := requestCount.Add(1)
				requestReachedServer <- struct{}{}
				if requestNumber == 1 {
					<-firstRequestMayComplete
				}
				w.Header().Set(HeaderCLIBackendStatus, "200")
				_, _ = w.Write([]byte("{}"))
			}))
			defer server.Close()

			serverURL, err := url.Parse(server.URL)
			assert.NoError(t, err)
			client := NewV2ClientWithHttpClient(server.Client(), serverURL, &RetryConfig{Enabled: false}, tt.options...)
			assert.Equal(t, tt.wantConcurrentRequests, client.allowConcurrentSessionRequests)
			client.session.Store(&Session{
				GlobalAccountSubdomain: "global-account",
				SessionId:              "session-id",
				IdentityProvider:       "custom-idp",
			})

			var requests sync.WaitGroup
			requestErrors := make(chan error, 2)
			requests.Add(2)
			for range 2 {
				go func() {
					defer requests.Done()
					_, err := client.Execute(context.Background(), NewGetRequest("test", map[string]string{}))
					requestErrors <- err
				}()
			}

			select {
			case <-requestReachedServer:
			case <-time.After(time.Second):
				t.Fatal("first request did not reach the server")
			}

			if tt.wantConcurrentRequests {
				select {
				case <-requestReachedServer:
				case <-time.After(time.Second):
					t.Fatal("second concurrent request did not reach the server")
				}
			} else {
				select {
				case <-requestReachedServer:
					t.Fatal("second request reached the server before the first completed")
				case <-time.After(100 * time.Millisecond):
				}
			}

			close(firstRequestMayComplete)
			requests.Wait()
			for range 2 {
				assert.NoError(t, <-requestErrors)
			}
			for range 2 {
				headers := <-requestHeaders
				assert.Equal(t, "session-id", headers.Get(HeaderCLISessionId))
				assert.Equal(t, "global-account", headers.Get(HeaderCLISubdomain))
				assert.Equal(t, "custom-idp", headers.Get(HeaderCLICustomIDP))
			}
		})
	}
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
				srvReturnContent: `{"issuer": "accounts.sap.com","user":"john.doe","mail":"john.doe@test.com"}`,
				srvReturnHeader: map[string]string{
					HeaderCLISessionId: "sessionid",
				},
				expectResponse: &LoginResponse{
					Issuer: "accounts.sap.com",
					Email:  "john.doe@test.com",
				},
				expectClientSession: &Session{
					SessionId:              "sessionid",
					IdentityProvider:       "",
					GlobalAccountSubdomain: "subdomain",
					LoggedInUser: &v2LoggedInUser{
						Issuer: "accounts.sap.com",
						Email:  "john.doe@test.com",
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
				srvReturnContent: `{"issuer": "customidp.accounts.ondemand.com","user":"john.doe","mail":"john.doe@test.com"}`,
				srvReturnHeader: map[string]string{
					HeaderCLISessionId: "sessionid",
				},
				expectResponse: &LoginResponse{
					Issuer: "customidp.accounts.ondemand.com",
					Email:  "john.doe@test.com",
				},
				expectClientSession: &Session{
					SessionId:              "sessionid",
					GlobalAccountSubdomain: "subdomain",
					IdentityProvider:       "my.custom.idp",
					LoggedInUser: &v2LoggedInUser{
						Issuer: "customidp.accounts.ondemand.com",
						Email:  "john.doe@test.com",
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
			description:  "error path - outdated protocol version [412]",
			loginRequest: NewLoginRequest("subdomain", "john.doe", "pass"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusPreconditionFailed,
				expectErrorMsg:  "Login failed due to outdated provider version. Update to the latest version of the provider. [Status: 412; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - unexpected error",
			loginRequest: NewLoginRequest("subdomain", "john.doe", "pass"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusTeapot,
				expectErrorMsg:  "received response with unexpected status: 418 [Status: 418; Correlation ID: fake-correlation-id]",
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

func TestV2Client_IdTokenLogin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		description string

		loginRequest *IdTokenLoginRequest
		simulation   v2SimulationConfig
	}{
		{
			description:  "happy path",
			loginRequest: NewIdTokenLoginRequest("subdomain", "idToken"),
			simulation: v2SimulationConfig{
				srvExpectBody:    `{"subdomain":"subdomain","idToken":"idToken"}`,
				srvReturnStatus:  http.StatusOK,
				srvReturnContent: `{"issuer": "idp.from.idtoken","user":"john.doe","mail":"john.doe@test.com"}`,
				srvReturnHeader: map[string]string{
					HeaderCLISessionId: "sessionid",
				},
				expectResponse: &LoginResponse{
					Issuer: "idp.from.idtoken",
					Email:  "john.doe@test.com",
				},
				expectClientSession: &Session{
					SessionId:              "sessionid",
					IdentityProvider:       "idp.from.idtoken",
					GlobalAccountSubdomain: "subdomain",
					LoggedInUser: &v2LoggedInUser{
						Issuer: "idp.from.idtoken",
						Email:  "john.doe@test.com",
					},
				},
			},
		},
		{
			description:  "error path - id token not parseable [400]",
			loginRequest: NewIdTokenLoginRequest("subdomain", "idToken"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusBadRequest,
				expectErrorMsg:  "Login failed. Invalid provider configuration. [Status: 400; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - id token expired [401]",
			loginRequest: NewIdTokenLoginRequest("subdomain", "idToken"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusUnauthorized,
				expectErrorMsg:  "Login failed. Please check ID Token validity. [Status: 401; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - globalaccount not found [404]",
			loginRequest: NewIdTokenLoginRequest("subdomain", "idToken"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusNotFound,
				expectErrorMsg:  "Global account 'subdomain' not found. Try again and make sure to provide the global account's subdomain. [Status: 404; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - outdated protocol version [412]",
			loginRequest: NewIdTokenLoginRequest("subdomain", "idToken"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusPreconditionFailed,
				expectErrorMsg:  "Login failed due to outdated provider version. Update to the latest version of the provider. [Status: 412; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - unexpected error",
			loginRequest: NewIdTokenLoginRequest("subdomain", "idToken"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusTeapot,
				expectErrorMsg:  "received response with unexpected status: 418 [Status: 418; Correlation ID: fake-correlation-id]",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			test.simulation.srvExpectPath = path.Join("/login", cliTargetProtocolVersion, "idtoken")
			test.simulation.callFunctionUnderTest = func(ctx context.Context, uut *v2Client) (any, error) {
				return uut.IdTokenLogin(ctx, test.loginRequest)
			}

			simulateV2Call(t, test.simulation)
		})
	}
}

func TestV2Client_BrowserLogin(t *testing.T) {

	isRealBrowser = false
	t.Parallel()

	// Setup HTTPS server
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			_, _ = w.Write([]byte("{\"loginId\":\"login.id\",\"subdomainRequired\":false}"))
		} else if strings.Contains(r.URL.Path, "/login.id") {

			res, _ := io.ReadAll(r.Body)

			idp := strings.Split(strings.Split(string(res), ",")[0], ":")[1]
			if idp == "\"\"" {
				idp = "accounts.sap.com"
			} else {
				idp = "customidp.accounts.ondemand.com"
			}

			w.Header().Add("X-Cpcli-Sessionid", "sessionid")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write([]byte("{\"issuer\":\"" + idp + "\",\"refreshToken\":\"sessionid\",\"user\":\"john.doe\",\"mail\":\"john.doe@test.com\"}"))
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}))

	srv.StartTLS()
	defer srv.Close()

	srvUrl, _ := url.Parse(srv.URL)

	t.Run("happy path", func(t *testing.T) {
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		uut.session.Store(&Session{})
		uut.newCorrelationID = func() string {
			return "fake-correlation-id"
		}
		_, err := uut.BrowserLogin(context.TODO(), &BrowserLoginRequest{
			GlobalAccountSubdomain: "my-subdomain",
		})

		assert.NoError(t, err)
		assert.Equal(t, &Session{
			SessionId:              "sessionid",
			IdentityProvider:       "",
			GlobalAccountSubdomain: "my-subdomain",
			LoggedInUser: &v2LoggedInUser{
				Issuer: "accounts.sap.com",
				Email:  "john.doe@test.com",
			},
		}, uut.session.Load())
	})

	t.Run("happy path - with custom idp", func(t *testing.T) {
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		uut.session.Store(&Session{})
		uut.newCorrelationID = func() string {
			return "fake-correlation-id"
		}
		_, err := uut.BrowserLogin(context.TODO(), &BrowserLoginRequest{
			CustomIdp:              "my.custom.idp",
			GlobalAccountSubdomain: "my-subdomain",
		})

		assert.NoError(t, err)
		assert.Equal(t, &Session{
			SessionId:              "sessionid",
			IdentityProvider:       "my.custom.idp",
			GlobalAccountSubdomain: "my-subdomain",
			LoggedInUser: &v2LoggedInUser{
				Issuer: "customidp.accounts.ondemand.com",
				Email:  "john.doe@test.com",
			},
		}, uut.session.Load())
	})
}

func TestV2Client_PasscodeLogin(t *testing.T) {
	t.Parallel()

	// Generate CA, server and client certificates
	caKey, caCert, _ := generateKeyPair(nil)
	serverKey, serverCert, _ := generateKeyPair(caKey)
	clientKey, clientCert, _ := generateKeyPair(caKey)

	_, pemEncodedCACerts := pemEncodeKeyPair(caKey, caCert)
	pemEncodedClientKey, pemEncodedClientCert := pemEncodeKeyPair(clientKey, clientCert)

	// Create certificate pool
	certPool := x509.NewCertPool()
	certPool.AddCert(caCert)

	// Create TLS config
	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{serverCert.Raw},
				PrivateKey:  serverKey,
			},
		},
		ClientCAs: certPool,
	}

	// Setup HTTPS server
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/service/users/passcode" {
			_, _ = w.Write([]byte("{\"passcode\": \"this-is-a-onetime-passcode\"}"))
		} else if strings.HasPrefix(r.URL.Path, "/login/") {
			var loginReq LoginRequest

			err := json.NewDecoder(r.Body).Decode(&loginReq)

			if assert.NoError(t, err) {
				assert.Equal(t, "this-is-a-onetime-passcode", loginReq.Password)
			}

			_, _ = w.Write([]byte("{}"))
		} else {
			w.WriteHeader(http.StatusNotImplemented)
		}
	}))
	srv.TLS = &tlsConfig
	srv.StartTLS()
	defer srv.Close()

	srvUrl, _ := url.Parse(srv.URL)

	t.Run("happy path", func(t *testing.T) {
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		_, err := uut.PasscodeLogin(context.TODO(), &PasscodeLoginRequest{
			GlobalAccountSubdomain: "my-subdomain",
			Username:               "john.doe@test.com",
			IdentityProvider:       "my-custom-idp",
			IdentityProviderURL:    srv.URL,
			PEMEncodedCACerts:      pemEncodedCACerts,
			PEMEncodedPrivateKey:   pemEncodedClientKey,
			PEMEncodedCertificate:  pemEncodedClientCert,
		})

		assert.NoError(t, err)
	})
	t.Run("error path - requires certificate", func(t *testing.T) {
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		_, err := uut.PasscodeLogin(context.TODO(), &PasscodeLoginRequest{
			GlobalAccountSubdomain: "my-subdomain",
			Username:               "john.doe@test.com",
			IdentityProvider:       "my-custom-idp",
			IdentityProviderURL:    srv.URL,
		})

		assert.Error(t, err)
	})
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
			_, _ = fmt.Fprintf(w, "{}")
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		_, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.NoError(t, err)
	})
	t.Run("backend headers get passed through", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(HeaderCLIBackendStatus, fmt.Sprintf("%d", 201))
			w.Header().Set(HeaderCLIBackendMediaType, "backend/mediatype")
			_, _ = fmt.Fprintf(w, "{}")
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		res, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.NoError(t, err)
		assert.Equal(t, 201, res.StatusCode)
		assert.Equal(t, "backend/mediatype", res.ContentType)
	})
	t.Run("session headers must be set", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "session-id", r.Header.Get(HeaderCLISessionId))
			assert.Equal(t, "globalaccount-subdomain", r.Header.Get(HeaderCLISubdomain))
			assert.Equal(t, "my.custom.idp", r.Header.Get(HeaderCLICustomIDP))
			w.Header().Set(HeaderCLIBackendStatus, fmt.Sprintf("%d", 201))
			w.Header().Set(HeaderCLIBackendMediaType, "backend/mediatype")
			_, _ = fmt.Fprintf(w, "{}")
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		uut.session.Store(&Session{
			GlobalAccountSubdomain: "globalaccount-subdomain",
			SessionId:              "session-id",
			IdentityProvider:       "my.custom.idp",
			LoggedInUser: &v2LoggedInUser{
				Email:  "john.doe@int.test",
				Issuer: "customidp.accounts.ondemand.com",
			},
		})

		cmdRes, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.NoError(t, err)
		assert.Equal(t, 201, cmdRes.StatusCode)
	})
	t.Run("backend error handling", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "my.custom.idp", r.Header.Get(HeaderCLICustomIDP))
			w.Header().Set(HeaderCLIBackendStatus, fmt.Sprintf("%d", 500))
			w.Header().Set(HeaderCLIBackendMediaType, "backend/mediatype")
			_, _ = fmt.Fprintf(w, `{"error":"this is a backend error"}`)
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		uut.session.Store(&Session{
			GlobalAccountSubdomain: "globalaccount-subdomain",
			IdentityProvider:       "my.custom.idp",
			LoggedInUser: &v2LoggedInUser{
				Email:  "john.doe@int.test",
				Issuer: "customidp.accounts.ondemand.com",
			},
		})

		cmdRes, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.EqualError(t, err, "this is a backend error")
		assert.Equal(t, 500, cmdRes.StatusCode)
	})
	t.Run("backend error handling - incompatible error message", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "my.custom.idp", r.Header.Get(HeaderCLICustomIDP))
			w.Header().Set(HeaderCLIBackendStatus, fmt.Sprintf("%d", 500))
			w.Header().Set(HeaderCLIBackendMediaType, "backend/mediatype")
			_, _ = fmt.Fprintf(w, `this is a backend error`)
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
		uut.session.Store(&Session{
			GlobalAccountSubdomain: "globalaccount-subdomain",
			IdentityProvider:       "my.custom.idp",
			LoggedInUser: &v2LoggedInUser{
				Email:  "john.doe@int.test",
				Issuer: "customidp.accounts.ondemand.com",
			},
		})

		cmdRes, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.EqualError(t, err, "the backend responded with an unknown error: 500")
		assert.Equal(t, 500, cmdRes.StatusCode)
	})
}

func TestV2Client_Execute_Tracing(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderCLIBackendStatus, "200")
		_, _ = fmt.Fprintf(w, "{}")
	}))
	defer srv.Close()

	srvUrl, _ := url.Parse(srv.URL)

	newTracingClient := func() (*v2Client, *bytes.Buffer) {
		var logSink bytes.Buffer
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, &RetryConfig{Enabled: false})
		uut.session.Store(&Session{GlobalAccountSubdomain: "my-global-account"})
		uut.newCorrelationID = func() string { return "test-correlation-id" }
		return uut, &logSink
	}

	findLogEntry := func(t *testing.T, logSink *bytes.Buffer, msg string) map[string]any {
		t.Helper()
		entries, err := tflogtest.MultilineJSONDecode(logSink)
		assert.NoError(t, err)
		for _, entry := range entries {
			if entry["@message"] == msg {
				return entry
			}
		}
		return nil
	}

	t.Run("logs command metadata at DEBUG level", func(t *testing.T) {
		uut, logSink := newTracingClient()
		ctx := tflogtest.RootLogger(context.Background(), logSink)

		_, err := uut.Execute(ctx, NewGetRequest("subaccount/role", map[string]string{}))
		assert.NoError(t, err)

		entry := findLogEntry(t, logSink, "Executing BTP CLI command")
		if assert.NotNil(t, entry, "expected DEBUG log entry 'Executing BTP CLI command' not found") {
			assert.Equal(t, "subaccount/role", entry["btp_command"])
			assert.Equal(t, string(ActionGet), entry["btp_action"])
			assert.Equal(t, "test-correlation-id", entry["btp_correlation_id"])
			assert.Equal(t, "my-global-account", entry["btp_global_account"])
		}
	})

	t.Run("logs subaccount ID at DEBUG level when present in args", func(t *testing.T) {
		uut, logSink := newTracingClient()
		ctx := tflogtest.RootLogger(context.Background(), logSink)

		_, err := uut.Execute(ctx, NewGetRequest("subaccount/role", map[string]string{
			"subaccount": "my-subaccount-id",
		}))
		assert.NoError(t, err)

		entry := findLogEntry(t, logSink, "Executing BTP CLI command")
		if assert.NotNil(t, entry, "expected DEBUG log entry 'Executing BTP CLI command' not found") {
			assert.Equal(t, "my-subaccount-id", entry["btp_subaccount_id"])
			assert.Nil(t, entry["btp_directory_id"])
		}
	})

	t.Run("logs directory ID at DEBUG level when directoryID key present in args", func(t *testing.T) {
		uut, logSink := newTracingClient()
		ctx := tflogtest.RootLogger(context.Background(), logSink)

		_, err := uut.Execute(ctx, NewGetRequest("accounts/directory", map[string]string{
			"directoryID": "my-directory-id",
		}))
		assert.NoError(t, err)

		entry := findLogEntry(t, logSink, "Executing BTP CLI command")
		if assert.NotNil(t, entry, "expected DEBUG log entry 'Executing BTP CLI command' not found") {
			assert.Equal(t, "my-directory-id", entry["btp_directory_id"])
			assert.Nil(t, entry["btp_subaccount_id"])
		}
	})

	t.Run("logs directory ID at DEBUG level when directory key present in args", func(t *testing.T) {
		uut, logSink := newTracingClient()
		ctx := tflogtest.RootLogger(context.Background(), logSink)

		_, err := uut.Execute(ctx, NewGetRequest("accounts/entitlement", map[string]string{
			"directory": "my-directory-id",
		}))
		assert.NoError(t, err)

		entry := findLogEntry(t, logSink, "Executing BTP CLI command")
		if assert.NotNil(t, entry, "expected DEBUG log entry 'Executing BTP CLI command' not found") {
			assert.Equal(t, "my-directory-id", entry["btp_directory_id"])
		}
	})

	t.Run("does not log subaccount or directory ID when not present in args", func(t *testing.T) {
		uut, logSink := newTracingClient()
		ctx := tflogtest.RootLogger(context.Background(), logSink)

		_, err := uut.Execute(ctx, NewGetRequest("accounts/global-account", map[string]string{}))
		assert.NoError(t, err)

		entry := findLogEntry(t, logSink, "Executing BTP CLI command")
		if assert.NotNil(t, entry, "expected DEBUG log entry 'Executing BTP CLI command' not found") {
			assert.Nil(t, entry["btp_subaccount_id"])
			assert.Nil(t, entry["btp_directory_id"])
		}
	})
}

func TestV2Client_DoRequest_Tracing(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderCLIBackendStatus, "200")
		_, _ = fmt.Fprintf(w, "{}")
	}))
	defer srv.Close()

	srvUrl, _ := url.Parse(srv.URL)

	t.Run("logs HTTP request details at TRACE level", func(t *testing.T) {
		var logSink bytes.Buffer
		ctx := tflogtest.RootLogger(context.Background(), &logSink)

		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, &RetryConfig{Enabled: false})
		uut.session.Store(&Session{GlobalAccountSubdomain: "my-global-account"})
		uut.newCorrelationID = func() string { return "test-correlation-id" }

		_, err := uut.Execute(ctx, NewGetRequest("subaccount/role", map[string]string{}))
		assert.NoError(t, err)

		entries, err := tflogtest.MultilineJSONDecode(&logSink)
		assert.NoError(t, err)

		var traceEntry map[string]any
		for _, entry := range entries {
			if entry["@message"] == "Sending BTP CLI HTTP request" {
				traceEntry = entry
				break
			}
		}

		if assert.NotNil(t, traceEntry, "expected TRACE log entry 'Sending BTP CLI HTTP request' not found") {
			assert.Equal(t, "POST", traceEntry["http_method"])
			assert.Contains(t, traceEntry["http_url"], "command/")
			assert.Equal(t, "test-correlation-id", traceEntry["btp_correlation_id"])
			assert.Equal(t, "my-global-account", traceEntry["btp_global_account"])
		}
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

	// the header values the fake api server shall return
	srvReturnHeader map[string]string

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
			assert.Equal(t, "Terraform/x.x.x terraform-plugin-btp/y.y.y", r.Header.Get("User-Agent"))

			if len(config.srvExpectBody) > 0 {
				assert.Equal(t, config.srvExpectBody, strings.TrimSpace(string(b)))
			}

			for key, value := range config.srvReturnHeader {
				w.Header().Add(key, value)
			}
			w.WriteHeader(config.srvReturnStatus)
			_, _ = fmt.Fprintf(w, "%s", config.srvReturnContent)
		}
	}))
	defer srv.Close()

	srvUrl, _ := url.Parse(srv.URL)
	uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl, nil)
	uut.UserAgent = "Terraform/x.x.x terraform-plugin-btp/y.y.y"
	uut.session.Store(config.initSession)
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

	assert.Equal(t, config.expectClientSession, uut.session.Load())
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
			Email:  "john.doe@test.com",
			Issuer: "test",
		}

		uut := NewV2Client(nil)
		uut.session.Store(&Session{LoggedInUser: testUser})

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
		uut.session.Store(&Session{GlobalAccountSubdomain: "my-subdomain"})

		assert.Equal(t, "my-subdomain", uut.GetGlobalAccountSubdomain())
	})
}

func generateKeyPair(signingKey *rsa.PrivateKey) (*rsa.PrivateKey, *x509.Certificate, error) {
	// generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	if signingKey == nil {
		signingKey = privateKey
	}

	// generate self-signed certificate
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(42),
		Subject: pkix.Name{
			Organization: []string{"SAP SE"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(5, 5, 5),
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, signingKey)
	if err != nil {
		return nil, nil, err
	}

	cert, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, cert, nil
}

func pemEncodeKeyPair(key *rsa.PrivateKey, cert *x509.Certificate) (pemEncodedKey string, pemEncodedCert string) {
	pemEncodedCert = string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw}))
	pemEncodedKey = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}))

	return
}

func TestV2Client_RetryLogic(t *testing.T) {
	t.Parallel()
	var attemptCount atomic.Int32
	retryStatus := 429

	// Simulate two rate limit failures, then allow success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentCount := attemptCount.Add(1)
		if currentCount <= 2 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(retryStatus)
			_, _ = w.Write([]byte(`rate limit`))
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"issuer": "accounts.sap.com","user":"john.doe","mail":"john.doe@test.com"}`))
		}
	}))
	defer server.Close()

	serverUrl, _ := url.Parse(server.URL)
	testClient := NewV2ClientWithHttpClient(server.Client(), serverUrl, nil)
	testClient.UserAgent = "TestAgent"
	testClient.newCorrelationID = func() string { return "test-cid" }

	// Minimal LoginRequest for trigger
	loginReq := &LoginRequest{
		GlobalAccountSubdomain: "subdomain",
		Username:               "john.doe",
		Password:               "pass",
	}

	start := time.Now()
	resp, err := testClient.Login(context.TODO(), loginReq)
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	// There should have been 3 attempts: 2 failures + 1 success
	assert.Equal(t, int32(3), attemptCount.Load())
	// The elapsed time should roughly reflect retries, but allow tolerance
	elapsedSec := elapsed.Seconds()
	assert.GreaterOrEqual(t, elapsedSec, 1.0, "should wait at least 1s due to Retry-After header")
	assert.Less(t, elapsedSec, 5.0, "should not take more than 5s even with retries")
}
