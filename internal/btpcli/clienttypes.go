package btpcli

import (
	"io"
)

/* Login */

// NewLoginRequest returns a new LoginRequest with `ldap` as default IdentityProvider set.
func NewLoginRequest(globalaccountSubdomain string, username string, password string) *LoginRequest {
	return NewLoginRequestWithCustomIDP("", globalaccountSubdomain, username, password)
}

func NewLoginRequestWithCustomIDP(idp string, globalaccountSubdomain string, username string, password string) *LoginRequest {
	return &LoginRequest{
		IdentityProvider:       idp,
		GlobalAccountSubdomain: globalaccountSubdomain,
		Username:               username,
		Password:               password,
	}
}

func NewIdTokenLoginRequest(globalaccountSubdomain string, idToken string) *IdTokenLoginRequest {
	return &IdTokenLoginRequest{
		GlobalAccountSubdomain: globalaccountSubdomain,
		IdToken:                idToken,
	}
}

type LoginRequest struct {
	IdentityProvider       string `json:"customIdp"`
	GlobalAccountSubdomain string `json:"subdomain"`
	Username               string `json:"userName"`
	Password               string `json:"password"`
}

type IdTokenLoginRequest struct {
	GlobalAccountSubdomain string `json:"subdomain"`
	IdToken                string `json:"idToken"`
}

type LoginResponse struct {
	RefreshToken string `json:"refreshToken"`
	Username     string `json:"user"`
	Email        string `json:"mail"`
	Issuer       string `json:"issuer"`
}

/* Logout */

func NewLogoutRequest(globalaccountSubdomain string) *LogoutRequest {
	return NewLogoutRequestWithCustomIDP("", globalaccountSubdomain)
}

func NewLogoutRequestWithCustomIDP(idp string, globalaccountSubdomain string) *LogoutRequest {
	return &LogoutRequest{
		IdentityProvider:       idp,
		GlobalAccountSubdomain: globalaccountSubdomain,
	}
}

type LogoutRequest struct {
	IdentityProvider       string `json:"customIdp"`
	GlobalAccountSubdomain string `json:"subdomain"`
	RefreshToken           string `json:"refreshToken"`
}

type LogoutResponse struct {
}

/* Command */
func NewCommandRequest(action Action, command string, args any) *CommandRequest {
	return &CommandRequest{
		Action:  action,
		Command: command,
		Args:    args,
	}
}

type CommandOptions struct {
	GoodState        int
	KnownErrorStates map[int]string
}

type Action string

type CommandRequest struct {
	Action  Action
	Command string
	Args    any
}

type CommandResponse struct {
	StatusCode  int
	ContentType string
	Body        io.ReadCloser
}
