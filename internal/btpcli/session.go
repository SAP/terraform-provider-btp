package btpcli

import "sync"

type v2LoggedInUser struct {
	Email  string
	Issuer string
}

type Session struct {
	GlobalAccountSubdomain string
	SessionId              string
	IdentityProvider       string
	LoggedInUser           *v2LoggedInUser

	// Mutex serializes requests that use this immutable session.
	sync.Mutex
}
