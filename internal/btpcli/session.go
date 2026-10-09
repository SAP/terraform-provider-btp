package btpcli

import (
	"sync"
)

type v2LoggedInUser struct {
	Email  string
	Issuer string
}

type Session struct {
	GlobalAccountSubdomain string
	SessionId              string
	IdentityProvider       string
	LoggedInUser           *v2LoggedInUser

	sync.RWMutex
}

type sessionSnapshot struct {
	globalAccountSubdomain string
	sessionId              string
	identityProvider       string
	loggedInUser           *v2LoggedInUser
}

func (s *Session) snapshot() sessionSnapshot {
	s.RLock()
	defer s.RUnlock()

	return sessionSnapshot{
		globalAccountSubdomain: s.GlobalAccountSubdomain,
		sessionId:              s.SessionId,
		identityProvider:       s.IdentityProvider,
		loggedInUser:           s.LoggedInUser,
	}
}
