package btpcli

import (
	"sync"
)

type v2LoggedInUser struct {
	Username string
	Email    string
	Issuer   string
}

type Session struct {
	GlobalAccountSubdomain string
	RefreshToken           string
	IdentityProvider       string
	LoggedInUser           *v2LoggedInUser

	sync.Mutex
}
