package btpcli

import (
	"context"
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

	sessionMutex
}

// sessionMutex serializes requests sharing a CLI session. Its zero value is ready
// for use, including sessions constructed directly by callers.
type sessionMutex struct {
	init  sync.Once
	token chan struct{}
}

func (m *sessionMutex) initialize() {
	m.init.Do(func() { m.token = make(chan struct{}, 1) })
}

func (m *sessionMutex) Lock() {
	_ = m.LockContext(context.Background())
}

func (m *sessionMutex) LockContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.initialize()
	select {
	case m.token <- struct{}{}:
		// Cancellation and acquisition may both be ready. Do not hand an expired
		// request the session in that case.
		if err := ctx.Err(); err != nil {
			m.Unlock()
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *sessionMutex) Unlock() {
	m.initialize()
	select {
	case <-m.token:
	default:
		panic("btpcli: unlock of unlocked session mutex")
	}
}
