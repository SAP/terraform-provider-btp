package btpcli

import (
	"sync"
)

type v2LoggedInUser struct {
	Username         string
	Email            string
	IdentityProvider string
}

type Session struct {
	GlobalAccountSubdomain string
	RefreshToken           string
	LoggedInUser           *v2LoggedInUser

	sync.Mutex
}
