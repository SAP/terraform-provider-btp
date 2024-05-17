package btpcli

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"

	uuid "github.com/hashicorp/go-uuid"
)

const DefaultServerURL string = "https://cli.btp.cloud.sap"

// We define an Uber Error type that is used to handle errors from the BTP CLI client.
// The error structure comprises the possible JSON structure of the error responses
type BtpClientError struct {
	Message     string         `json:"error"` // This should always be available
	Description string         `json:"description"`
	BrokerError *SmBrokerError `json:"broker_error"` // Service Broker/Service Manager specific
}

// The service broker error
type SmBrokerError struct {
	// The broker status code.
	StatusCode int32 `json:"StatusCode"`
	// A machine-readable error string that may be returned by the broker.
	ErrorMessage string `json:"ErrorMessage"`
	// A human-readable description of the error that may be returned by the broker.
	Description string `json:"Description"`
	// ResponseError is set to the error that occurred when unmarshalling a response body from the broker.
	ResponseError string `json:"ResponseError"`
}

func NewV2Client(serverURL *url.URL) *v2Client {
	return NewV2ClientWithHttpClient(http.DefaultClient, serverURL)
}

func NewV2ClientWithHttpClient(client *http.Client, serverURL *url.URL) *v2Client {
	return &v2Client{
		httpClient: injectBTPCLITransport(client),
		serverURL:  serverURL,
		newCorrelationID: func() string {
			val, err := uuid.GenerateUUID()
			if err != nil {
				panic(fmt.Sprintf("crypto/rand returned fatal error: %s", err.Error()))
			}
			return val
		},
	}
}

const (
	HeaderCorrelationID       string = "X-Correlationid"
	HeaderIDToken             string = "X-Id-Token"
	HeaderCLIFormat           string = "X-Cpcli-Format"
	HeaderCLISessionId        string = "X-Cpcli-Sessionid"
	HeaderCLISubdomain        string = "X-Cpcli-Subdomain"
	HeaderCLICustomIDP        string = "X-Cpcli-Customidp"
	HeaderCLIBackendStatus    string = "X-Cpcli-Backend-Status"
	HeaderCLIBackendMessage   string = "X-Cpcli-Backend-Message"
	HeaderCLIBackendMediaType string = "X-Cpcli-Backend-Mediatype"
	HeaderCLIClientUpdate     string = "X-Cpcli-Client-Update"
	HeaderCLIServerMessage    string = "X-Cpcli-Server-Message"
)

const cliTargetProtocolVersion string = "v2.49.0"

type v2ContextKey string

type v2Client struct {
	httpClient *http.Client
	serverURL  *url.URL

	newCorrelationID func() string

	session   *Session
	UserAgent string
}

func (v2 *v2Client) initTrace(ctx context.Context) context.Context {
	return context.WithValue(ctx, v2ContextKey(HeaderCorrelationID), v2.newCorrelationID())
}

func (v2 *v2Client) doRequest(ctx context.Context, method string, endpoint string, body any) (*http.Response, error) {
	endpointURL, err := url.Parse(endpoint)

	if err != nil {
		return nil, err
	}

	var bodyContent bytes.Buffer

	if body != nil {
		encoder := json.NewEncoder(&bodyContent)
		encoder.SetEscapeHTML(false) // no need for HTML safe json encoding
		if err = encoder.Encode(body); err != nil {
			return nil, err
		}
	}

	fullQualifiedEndpointURL := v2.serverURL.ResolveReference(endpointURL)

	req, err := http.NewRequestWithContext(ctx, method, fullQualifiedEndpointURL.String(), &bodyContent)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", v2.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderCLIFormat, "json")

	if v2.session != nil {
		v2.session.Lock()
		defer v2.session.Unlock()

		req.Header.Set(HeaderCLISessionId, v2.session.SessionId)
		req.Header.Set(HeaderCLISubdomain, v2.session.GlobalAccountSubdomain)
		req.Header.Set(HeaderCLICustomIDP, v2.session.IdentityProvider)
	}

	if correlationID := ctx.Value(v2ContextKey(HeaderCorrelationID)); correlationID != nil {
		req.Header.Set(HeaderCorrelationID, correlationID.(string))
	}

	res, err := v2.httpClient.Do(req)

	return res, err
}

func (v2 *v2Client) doPostRequest(ctx context.Context, endpoint string, body any) (*http.Response, error) {
	return v2.doRequest(ctx, http.MethodPost, endpoint, body)
}

func (v2 *v2Client) doGetRequest(ctx context.Context, endpoint string) (*http.Response, error) {
	return v2.doRequest(ctx, http.MethodGet, endpoint, nil)
}

func (v2 *v2Client) parseResponse(ctx context.Context, res *http.Response, targetObj any, goodState int, knownErrorStates map[int]string) error {
	if err := v2.checkResponseForErrors(ctx, res, goodState, knownErrorStates); err != nil {
		return err
	}

	if targetObj == nil {
		return nil
	}

	return json.NewDecoder(res.Body).Decode(targetObj)
}

func (v2 *v2Client) checkResponseForErrors(ctx context.Context, res *http.Response, goodState int, knownErrorStates map[int]string) error {
	if res.StatusCode == goodState {
		return nil
	}

	var err error

	if errorMsg, known := knownErrorStates[res.StatusCode]; known {
		err = fmt.Errorf("%s", errorMsg)
	} else {
		err = v2.parseResponseError(ctx, res)
	}

	return fmt.Errorf("%w [Status: %d; Correlation ID: %s]", err, res.StatusCode, ctx.Value(v2ContextKey(HeaderCorrelationID)))
}

func (v2 *v2Client) parseResponseError(ctx context.Context, res *http.Response) error {
	return fmt.Errorf("received response with unexpected status")
}

// Login authenticates a user using username + password
func (v2 *v2Client) Login(ctx context.Context, loginReq *LoginRequest) (*LoginResponse, error) {
	ctx = v2.initTrace(ctx)

	// TODO: After the switch to client protocol v2.49.0 the terraform provider is still providing
	//       the globalaccount subdomain during login. However, this relies on a special handling
	//       for older clients that might be removed from the server in the future.
	res, err := v2.doPostRequest(ctx, path.Join("login", cliTargetProtocolVersion), loginReq)

	if err != nil {
		return nil, err
	}

	var loginResponse LoginResponse
	err = v2.parseResponse(ctx, res, &loginResponse, http.StatusOK, map[int]string{
		http.StatusUnauthorized:       "Login failed. Check your credentials.",
		http.StatusForbidden:          fmt.Sprintf("You cannot access global account '%s'. Make sure you have at least read access to the global account, a directory, or a subaccount.", loginReq.GlobalAccountSubdomain),
		http.StatusNotFound:           fmt.Sprintf("Global account '%s' not found. Try again and make sure to provide the global account's subdomain.", loginReq.GlobalAccountSubdomain),
		http.StatusPreconditionFailed: "Login failed due to outdated provider version. Update to the latest version of the provider.",
		http.StatusGatewayTimeout:     "Login timed out. Please try again later.",
	})

	if err != nil {
		return nil, err
	}

	v2.session = &Session{
		GlobalAccountSubdomain: loginReq.GlobalAccountSubdomain,
		IdentityProvider:       loginReq.IdentityProvider,
		LoggedInUser: &v2LoggedInUser{
			Username: loginResponse.Username,
			Email:    loginResponse.Email,
			Issuer:   loginResponse.Issuer,
		},
		SessionId: res.Header.Get(HeaderCLISessionId),
	}

	return &loginResponse, nil
}

// IdTokenLogin authenticates a user by providing an id token
func (v2 *v2Client) IdTokenLogin(ctx context.Context, loginReq *IdTokenLoginRequest) (*LoginResponse, error) {
	ctx = v2.initTrace(ctx)

	res, err := v2.doPostRequest(ctx, path.Join("login", cliTargetProtocolVersion, "idtoken"), loginReq)

	if err != nil {
		return nil, err
	}

	var loginResponse LoginResponse
	err = v2.parseResponse(ctx, res, &loginResponse, http.StatusOK, map[int]string{
		http.StatusBadRequest:         "Login failed. Invalid provider configuration.",
		http.StatusUnauthorized:       "Login failed. Please check ID Token validity.",
		http.StatusNotFound:           fmt.Sprintf("Global account '%s' not found. Try again and make sure to provide the global account's subdomain.", loginReq.GlobalAccountSubdomain),
		http.StatusPreconditionFailed: "Login failed due to outdated provider version. Update to the latest version of the provider.",
		http.StatusGatewayTimeout:     "Login timed out. Please try again later.",
	})

	if err != nil && err.Error() != "EOF" { // TODO: stop ignoring EOF when btp CLI server returning non-empty body has reached productive landscapes
		return nil, err
	}

	v2.session = &Session{
		GlobalAccountSubdomain: loginReq.GlobalAccountSubdomain,
		IdentityProvider:       loginResponse.Issuer,
		LoggedInUser: &v2LoggedInUser{
			Username: loginResponse.Username,
			Email:    loginResponse.Email,
			Issuer:   loginResponse.Issuer,
		},
		SessionId: res.Header.Get(HeaderCLISessionId),
	}

	return &loginResponse, nil
}

// BrowserLogin authenticates user using SSO
func (v2 *v2Client) BrowserLogin(ctx context.Context, loginReq *BrowserLoginRequest) (*BrowserLoginPostResponse, error) {
	ctx = v2.initTrace(ctx)

	// TODO: After the switch to client protocol v2.49.0 the terraform provider is still providing
	//       the globalaccount subdomain during login. However, this relies on a special handling
	//       for older clients that might be removed from the server in the future.

	optionalIdpPath := OptionalCustomIdpPath(loginReq)
	res, err := v2.doGetRequest(ctx, path.Join("login", cliTargetProtocolVersion, "browser", optionalIdpPath))

	if err != nil {
		return nil, err
	}

	var browserResponse BrowserResponse
	err = v2.parseResponse(ctx, res, &browserResponse, http.StatusOK, map[int]string{
		http.StatusUnauthorized:       "Login failed. Check your credentials.",
		http.StatusPreconditionFailed: "Login failed due to outdated provider version. Update to the latest version of the provider.",
		http.StatusGatewayTimeout:     "Login timed out. Please try again later.",
	})

	if err != nil {
		return nil, err
	}

	endpointURL, err := url.Parse(path.Join("login", cliTargetProtocolVersion, "browser", browserResponse.LoginID, optionalIdpPath))
	if err != nil {
		return nil, err
	}

	fullQualifiedEndpointURL := v2.serverURL.ResolveReference(endpointURL)
	err = openUserAgent(fullQualifiedEndpointURL.String(), isRealBrowser)

	if err != nil {
		fmt.Printf("browser_login_open_browser_failed : %s", fullQualifiedEndpointURL.String())
		return nil, err
	} else {
		GiveBrowserTimeToOpen()
	}

	res, err = v2.doPostRequest(ctx, path.Join("login", cliTargetProtocolVersion, "browser", browserResponse.LoginID), loginReq)
	if err != nil {
		return nil, err
	}

	var browserLoginPostResponse BrowserLoginPostResponse
	err = v2.parseResponse(ctx, res, &browserLoginPostResponse, http.StatusOK, map[int]string{
		http.StatusUnauthorized:       "Login failed. Check your credentials.",
		http.StatusForbidden:          fmt.Sprintf("You cannot access global account '%s'. Make sure you have at least read access to the global account, a directory, or a subaccount.", loginReq.GlobalAccountSubdomain),
		http.StatusNotFound:           fmt.Sprintf("Global account '%s' not found. Try again and make sure to provide the global account's subdomain.", loginReq.GlobalAccountSubdomain),
		http.StatusPreconditionFailed: "Login failed due to outdated provider version. Update to the latest version of the provider.",
		http.StatusGatewayTimeout:     "Login timed out. Please try again later.",
	})
	if err != nil {
		return nil, err
	}

	v2.session = &Session{
		GlobalAccountSubdomain: loginReq.GlobalAccountSubdomain,
		IdentityProvider:       loginReq.CustomIdp,
		LoggedInUser: &v2LoggedInUser{
			Username: browserLoginPostResponse.Username,
			Email:    browserLoginPostResponse.Email,
			Issuer:   browserLoginPostResponse.Issuer,
		},
		SessionId: browserLoginPostResponse.RefreshToken,
	}

	return &browserLoginPostResponse, nil
}

/*
The variable isRealBrowser primarily serves testing purposes. During the testing of the browser login flow, the  variable is intentionally set to
false. This configuration is allows the stubbing of the call that opens the browser, as no validation is necessary in this scenario.
*/
var isRealBrowser = true

func openUserAgent(url string, isRealBrowser bool) error {
	if isRealBrowser {
		switch runtime.GOOS {
		case "linux":
			return exec.Command("xdg-open", url).Start()
		case "windows":
			return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			return exec.Command("open", url).Start()
		default:
			return fmt.Errorf("unsupported_platform")
		}
	} else {
		return nil
	}
}

func GiveBrowserTimeToOpen() {
	time.Sleep(1 * time.Second)
}

func OptionalCustomIdpPath(loginReq *BrowserLoginRequest) string {
	if loginReq.CustomIdp == "" {
		return ""
	}
	return "/idp/" + loginReq.CustomIdp
}

// PasscodeLogin authenticates with a pem encoded x509 key-pair
func (v2 *v2Client) PasscodeLogin(ctx context.Context, loginReq *PasscodeLoginRequest) (*LoginResponse, error) {
	ctx = v2.initTrace(ctx)

	clientCert, err := tls.X509KeyPair([]byte(loginReq.PEMEncodedCertificate), []byte(loginReq.PEMEncodedPrivateKey))
	if err != nil {
		return nil, err
	}

	caCertPool, err := x509.SystemCertPool()

	if err != nil {
		return nil, err
	}

	if len(loginReq.PEMEncodedCACerts) > 0 {
		caCertPool.AppendCertsFromPEM([]byte(loginReq.PEMEncodedCACerts))
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	tlsTransport := (http.DefaultTransport.(*http.Transport)).Clone()
	tlsTransport.TLSClientConfig = tlsConfig
	idpClient := &http.Client{Transport: tlsTransport}

	res, err := idpClient.Get(fmt.Sprintf("%s/service/users/passcode", loginReq.IdentityProviderURL)) // TODO use URL

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IDP responded with unexpected response code: %d", res.StatusCode)
	}

	var passcodeResponse struct {
		Passcode string `json:"passcode"`
	}

	if err := json.NewDecoder(res.Body).Decode(&passcodeResponse); err != nil {
		return nil, err
	}

	return v2.Login(ctx, &LoginRequest{
		IdentityProvider:       loginReq.IdentityProvider,
		GlobalAccountSubdomain: loginReq.GlobalAccountSubdomain,
		Username:               loginReq.Username,
		Password:               passcodeResponse.Passcode,
	})
}

// Execute executes a command
func (v2 *v2Client) Execute(ctx context.Context, cmdReq *CommandRequest, options ...CommandOptions) (cmdRes CommandResponse, err error) {
	ctx = v2.initTrace(ctx)

	wrappedArgs := struct {
		ParamValues any `json:"paramValues"`
	}{
		ParamValues: cmdReq.Args,
	}

	res, err := v2.doPostRequest(ctx, fmt.Sprintf("%s?%s", path.Join("command", cliTargetProtocolVersion, cmdReq.Command), cmdReq.Action), wrappedArgs)

	if err != nil {
		return
	}

	opts := firstElementOrDefault(options, CommandOptions{GoodState: http.StatusOK, KnownErrorStates: map[int]string{}})
	opts.KnownErrorStates[http.StatusGatewayTimeout] = "Command timed out. Please try again later."
	opts.KnownErrorStates[http.StatusForbidden] = "Access forbidden due to insufficient authorization. Make sure to have sufficient access rights."

	if err = v2.checkResponseForErrors(ctx, res, opts.GoodState, opts.KnownErrorStates); err != nil {
		return
	}

	if cmdRes.StatusCode, err = strconv.Atoi(res.Header.Get(HeaderCLIBackendStatus)); err != nil {
		err = fmt.Errorf("unable to convert reported backend status code: %w", err)
		return
	}

	if cmdRes.StatusCode >= 400 {

		var backendError BtpClientError

		if err = json.NewDecoder(res.Body).Decode(&backendError); err == nil {
			err = fmt.Errorf(backendError.Message)
			// Handle scenarios where more error context can potentially be provided
			err = handleSpecialErrors(backendError, err)
		} else if res.Header.Get(HeaderCLIServerMessage) != "" {
			err = fmt.Errorf("the backend responded with an error: %s", res.Header.Get(HeaderCLIServerMessage))
		} else {
			err = fmt.Errorf("the backend responded with an unknown error: %d", cmdRes.StatusCode)
		}

		return
	}

	cmdRes.Body = res.Body
	cmdRes.ContentType = res.Header.Get(HeaderCLIBackendMediaType)
	return
}

func handleSpecialErrors(backendError BtpClientError, plainError error) error {
	// Errors that go beyond the plain error message can be handled in this function
	if backendError.BrokerError != nil {
		// Handle the additional information provided by the service broker/service manager
		return fmt.Errorf("%s - %s", backendError.Message, backendError.BrokerError.Description)
	}

	return plainError
}

func (v2 *v2Client) GetGlobalAccountSubdomain() string {
	if v2.session == nil {
		return ""
	}

	return v2.session.GlobalAccountSubdomain
}

func (v2 *v2Client) GetLoggedInUser() *v2LoggedInUser {
	if v2.session == nil {
		return nil
	}

	return v2.session.LoggedInUser
}
