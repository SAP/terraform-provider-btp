package btpcli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	uuid "github.com/hashicorp/go-uuid"
)

const DefaultServerURL string = "https://cpcli.cf.eu10.hana.ondemand.com"

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
	HeaderCorrelationID              string = "X-CorrelationId"
	HeaderIDToken                    string = "X-Id-Token"
	HeaderCLIFormat                  string = "X-CPCLI-Format"
	HeaderCLIRefreshToken            string = "X-CPCLI-RefreshToken"
	HeaderCLIReplacementRefreshToken string = "X-CPCLI-ReplacementRefreshtoken"
	HeaderCLISubdomain               string = "X-CPCLI-Subdomain"
	HeaderCLICustomIDP               string = "X-CPCLI-CustomIdp"
	HeaderCLIBackendStatus           string = "X-CPCLI-Backend-Status"
	HeaderCLIBackendMessage          string = "X-CPCLI-Backend-Message"
	HeaderCLIBackendMediaType        string = "X-CPCLI-Backend-MediaType"
)

const cliTargetProtocolVersion string = "v2.38.0"

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
		if err = json.NewEncoder(&bodyContent).Encode(body); err != nil {
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

		req.Header.Set(HeaderCLIRefreshToken, v2.session.RefreshToken)
		req.Header.Set(HeaderCLISubdomain, v2.session.GlobalAccountSubdomain)
		req.Header.Set(HeaderCLICustomIDP, v2.session.IdentityProvider)
	}

	if correlationID := ctx.Value(v2ContextKey(HeaderCorrelationID)); correlationID != nil {
		req.Header.Set(HeaderCorrelationID, correlationID.(string))
	}

	res, err := v2.httpClient.Do(req)

	if v2.session != nil && err == nil {
		v2.session.RefreshToken = res.Header.Get(HeaderCLIReplacementRefreshToken)
	}

	return res, err
}

func (v2 *v2Client) doPostRequest(ctx context.Context, endpoint string, body any) (*http.Response, error) {
	return v2.doRequest(ctx, http.MethodPost, endpoint, body)
}

func (v2 *v2Client) parseResponse(ctx context.Context, res *http.Response, targetObj any, goodState int, knownErrorStates map[int]string) error {
	if err := v2.checkResponseForErrors(ctx, res, goodState, knownErrorStates); err != nil {
		return err
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
	return fmt.Errorf("Received response with unexpected status")
}

// Login authenticates a user using username + password
func (v2 *v2Client) Login(ctx context.Context, loginReq *LoginRequest) (*LoginResponse, error) {
	ctx = v2.initTrace(ctx)

	res, err := v2.doPostRequest(ctx, path.Join("login", cliTargetProtocolVersion), loginReq)

	if err != nil {
		return nil, err
	}

	var loginResponse LoginResponse
	err = v2.parseResponse(ctx, res, &loginResponse, http.StatusOK, map[int]string{
		http.StatusUnauthorized:   "Login failed. Check your credentials.",
		http.StatusForbidden:      fmt.Sprintf("You cannot access global account '%s'. Make sure you have at least read access to the global account, a directory, or a subaccount.", loginReq.GlobalAccountSubdomain),
		http.StatusNotFound:       fmt.Sprintf("Global account '%s' not found. Try again and make sure to provide the global account's subdomain.", loginReq.GlobalAccountSubdomain),
		http.StatusGatewayTimeout: "Login timed out. Please try again later.",
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
		RefreshToken: loginResponse.RefreshToken,
	}

	return &loginResponse, nil
}

// Logout invalidates the current user session
func (v2 *v2Client) Logout(ctx context.Context, logoutReq *LogoutRequest) (*LogoutResponse, error) {
	ctx = v2.initTrace(ctx)

	res, err := v2.doPostRequest(ctx, path.Join("logout", cliTargetProtocolVersion), logoutReq)

	if err != nil {
		return nil, err
	}

	var logoutResponse LogoutResponse
	return &logoutResponse, v2.parseResponse(ctx, res, &logoutResponse, http.StatusOK, map[int]string{
		http.StatusGatewayTimeout: "Logout timed out. Please try again later.",
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

	if err = v2.checkResponseForErrors(ctx, res, opts.GoodState, opts.KnownErrorStates); err != nil {
		return
	}

	if cmdRes.StatusCode, err = strconv.Atoi(res.Header.Get(HeaderCLIBackendStatus)); err != nil {
		err = fmt.Errorf("unable to convert reported backend status code: %w", err)
		return
	}

	if cmdRes.StatusCode >= 400 {
		var backendError struct {
			Message string `json:"error"`
		}

		if err = json.NewDecoder(res.Body).Decode(&backendError); err == nil {
			err = fmt.Errorf(backendError.Message)
		} else {
			err = fmt.Errorf("the backend responded with an unknown error: %d", cmdRes.StatusCode)
		}

		return
	}

	cmdRes.Body = res.Body
	cmdRes.ContentType = res.Header.Get(HeaderCLIBackendMediaType)
	return
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
