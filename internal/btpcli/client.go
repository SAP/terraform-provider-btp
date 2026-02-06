package btpcli

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	uuid "github.com/hashicorp/go-uuid"
)

const DefaultServerURL string = "https://cli.btp.cloud.sap"

// We define an Uber Error type that is used to handle errors from the BTP CLI client.
// The error structure comprises the possible JSON structure of the error responses
// Some services return "error", others return "ErrorMessage". Both are mapped here for consistent handling.
type BtpClientError struct {
	Message      string         `json:"error"` // This should always be available
	Description  string         `json:"description"`
	BrokerError  *SmBrokerError `json:"broker_error"` // Service Broker/Service Manager specific
	DestError    *[]DestError   `json:"violations"`   // Destination service specific error details
	ErrorMessage string         `json:"ErrorMessage"` // Alternate error field used by some services
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

// Destination service specific error details
type DestError struct {
	// A list of specific error details that may be returned by the destination service.
	Configuration string   `json:"configuration"`
	Errors        []string `json:"errors"`
}

type ErrorResponseBody struct {
	Error *ResponseBodyError `json:"error"`
}

type ResponseBodyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Target  string `json:"target"`
	Details []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"details"`
}

type RetryConfig struct {
	Enabled      bool
	RetryMax     int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
}

func NewV2Client(serverURL *url.URL) *v2Client {
	return NewV2ClientWithHttpClient(http.DefaultClient, serverURL, nil)
}

func NewRetryableHttpClient(cfg *RetryConfig) *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	if cfg == nil {
		cfg = &RetryConfig{
			Enabled:      true,
			RetryMax:     6,
			RetryWaitMin: 1 * time.Second,
			RetryWaitMax: 120 * time.Second,
		}
	}
	retryClient.RetryMax = cfg.RetryMax
	retryClient.RetryWaitMin = cfg.RetryWaitMin
	retryClient.RetryWaitMax = cfg.RetryWaitMax
	retryClient.Logger = nil

	if !cfg.Enabled {
		retryClient.RetryMax = 0
		retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			return false, nil
		}
		retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
			return 0
		}
		return retryClient
	}

	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// Retry on transient network errors and specific HTTP status codes (429, 500, 502, 503)
		if err != nil {
			return true, nil
		}
		if resp == nil {
			return false, nil
		}

		switch resp.StatusCode {
		case http.StatusTooManyRequests, // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusServiceUnavailable,  // 503
			http.StatusGatewayTimeout:      // 504
			// retry only these
			return true, nil

		case http.StatusBadRequest: // 400
			// Peek into the body to check for specific error codes/messages
			const maxBodyPeek = 4096
			var buf bytes.Buffer
			tee := io.TeeReader(io.LimitReader(resp.Body, maxBodyPeek), &buf)
			peekBytes, _ := io.ReadAll(tee)

			resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(buf.Bytes()), resp.Body))

			if strings.Contains(string(peekBytes), "[Error: 30004/400]") {
				return true, nil // for locking scenario API call must be retried
			}
			return false, nil
		default:
			// do not retry on 4xx client errors, or other 5xx errors
			return false, nil
		}
	}
	return retryClient
}

func NewV2ClientWithHttpClient(client *http.Client, serverURL *url.URL, retryCfg *RetryConfig) *v2Client {
	retryClient := NewRetryableHttpClient(retryCfg)
	retryClient.HTTPClient = client
	return &v2Client{
		httpClient: injectBTPCLITransport(retryClient.StandardClient()),
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

const cliTargetProtocolVersion string = "v2.97.0"

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
		err = v2.parseResponseError(res)
	}

	return fmt.Errorf("%w [Status: %d; Correlation ID: %s]", err, res.StatusCode, ctx.Value(v2ContextKey(HeaderCorrelationID)))
}

func (v2 *v2Client) parseResponseError(res *http.Response) error {
	if res.StatusCode == 429 {
		var errorBody ErrorResponseBody
		// Decode error body JSON
		if err := json.NewDecoder(res.Body).Decode(&errorBody); err != nil {
			return fmt.Errorf("rate limit exceeded (HTTP 429). Failed to parse error details: %v", err)
		}

		if errorBody.Error == nil {
			return fmt.Errorf("rate limit exceeded (HTTP 429). Response body missing 'error' field")
		}

		if errorBody.Error.Code == 11006 {
			var detail []string
			for _, d := range errorBody.Error.Details {
				detail = append(detail, d.Message)
			}
			details := strings.Join(detail, "; ")
			if details != "" {
				details = ": " + details
			}

			return fmt.Errorf("rate limit exceeded (HTTP 429) for target '%s': %s%s", errorBody.Error.Target, errorBody.Error.Message, details)
		}

		// If code not 11006, fallback message
		return fmt.Errorf("received HTTP 429 but unexpected error code %d: %s", errorBody.Error.Code, errorBody.Error.Message)
	}

	return fmt.Errorf("received response with unexpected status: %d", res.StatusCode)
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
			Email:  loginResponse.Email,
			Issuer: loginResponse.Issuer,
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
			Email:  loginResponse.Email,
			Issuer: loginResponse.Issuer,
		},
		SessionId: res.Header.Get(HeaderCLISessionId),
	}

	return &loginResponse, nil
}

// BrowserLogin authenticates user using SSO
func (v2 *v2Client) BrowserLogin(ctx context.Context, loginReq *BrowserLoginRequest) (*LoginResponse, error) {
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

	var browserLoginPostResponse LoginResponse
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
			Email:  browserLoginPostResponse.Email,
			Issuer: browserLoginPostResponse.Issuer,
		},
		SessionId: res.Header.Get(HeaderCLISessionId),
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
			msg := backendError.Message
			if msg == "" && backendError.ErrorMessage != "" {
				msg = backendError.ErrorMessage
			}
			if msg == "" && backendError.Description != "" {
				msg = backendError.Description
			}
			if msg == "" {
				msg = "unknown backend error"
			}
			err = fmt.Errorf("%s", msg)
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

	if backendError.DestError != nil {
		// Handle the additional information provided by the destination service
		var errorDetails []string
		for _, violation := range *backendError.DestError {
			if len(violation.Errors) > 0 {
				configErrors := fmt.Sprintf("Configuration '%s': %s", violation.Configuration, strings.Join(violation.Errors, "; "))
				errorDetails = append(errorDetails, configErrors)
			}
		}
		if len(errorDetails) > 0 {
			return fmt.Errorf("%s - %s", backendError.ErrorMessage, strings.Join(errorDetails, " | "))
		}
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
