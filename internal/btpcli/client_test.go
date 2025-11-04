package btpcli

import (
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
	"testing"
	"time"

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
				expectErrorMsg:  "received response with unexpected status [Status: 418; Correlation ID: fake-correlation-id]",
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
			description:  "error path - login request times out [504]]",
			loginRequest: NewIdTokenLoginRequest("subdomain", "idToken"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusGatewayTimeout,
				expectErrorMsg:  "Login timed out. Please try again later. [Status: 504; Correlation ID: fake-correlation-id]",
			},
		},
		{
			description:  "error path - unexpected error",
			loginRequest: NewIdTokenLoginRequest("subdomain", "idToken"),
			simulation: v2SimulationConfig{
				srvReturnStatus: http.StatusTeapot,
				expectErrorMsg:  "received response with unexpected status [Status: 418; Correlation ID: fake-correlation-id]",
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
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
		uut.session = &Session{}
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
		}, uut.session)
	})

	t.Run("happy path - with custom idp", func(t *testing.T) {
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
		uut.session = &Session{}
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
		}, uut.session)
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
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
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
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
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
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
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
			_, _ = fmt.Fprintf(w, "{}")
		}))
		defer srv.Close()

		srvUrl, _ := url.Parse(srv.URL)
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
		uut.session = &Session{
			GlobalAccountSubdomain: "globalaccount-subdomain",
			IdentityProvider:       "my.custom.idp",
			LoggedInUser: &v2LoggedInUser{
				Email:  "john.doe@int.test",
				Issuer: "customidp.accounts.ondemand.com",
			},
		}

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
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
		uut.session = &Session{
			GlobalAccountSubdomain: "globalaccount-subdomain",
			IdentityProvider:       "my.custom.idp",
			LoggedInUser: &v2LoggedInUser{
				Email:  "john.doe@int.test",
				Issuer: "customidp.accounts.ondemand.com",
			},
		}

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
		uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
		uut.session = &Session{
			GlobalAccountSubdomain: "globalaccount-subdomain",
			IdentityProvider:       "my.custom.idp",
			LoggedInUser: &v2LoggedInUser{
				Email:  "john.doe@int.test",
				Issuer: "customidp.accounts.ondemand.com",
			},
		}

		cmdRes, err := uut.Execute(context.TODO(), NewGetRequest("subaccount/role", map[string]string{}))

		assert.EqualError(t, err, "the backend responded with an unknown error: 500")
		assert.Equal(t, 500, cmdRes.StatusCode)
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
	uut := NewV2ClientWithHttpClient(srv.Client(), srvUrl)
	uut.UserAgent = "Terraform/x.x.x terraform-plugin-btp/y.y.y"
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
			Email:  "john.doe@test.com",
			Issuer: "test",
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
	attemptCount := 0
	retryStatus := 429

	// Simulate two rate limit failures, then allow success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount <= 2 {
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
	testClient := NewV2ClientWithHttpClient(server.Client(), serverUrl)
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
	assert.Equal(t, 3, attemptCount)
	// The time elapsed should be at least 2 seconds due to Retry-After header
	assert.GreaterOrEqual(t, int(elapsed.Seconds()), 2)
}
