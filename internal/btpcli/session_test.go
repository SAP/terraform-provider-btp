package btpcli

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSessionLockContext(t *testing.T) {
	for _, deadline := range []bool{false, true} {
		t.Run(map[bool]string{false: "cancellation", true: "deadline"}[deadline], func(t *testing.T) {
			session := &Session{}
			session.Lock()
			defer session.Unlock()
			ctx, cancel := context.WithCancel(context.Background())
			expected := context.Canceled
			if deadline {
				cancel()
				ctx, cancel = context.WithTimeout(context.Background(), 20*time.Millisecond)
				expected = context.DeadlineExceeded
			}
			defer cancel()
			done := make(chan error, 1)
			go func() { done <- session.LockContext(ctx) }()
			if !deadline {
				cancel()
			}
			select {
			case err := <-done:
				if !errors.Is(err, expected) {
					t.Fatalf("expected %v, got %v", expected, err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("waiter did not return while the session remained held")
			}
		})
	}
}

func TestSessionLockContextAlreadyCanceled(t *testing.T) {
	session := &Session{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for range 100 {
		if err := session.LockContext(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	}
	if err := session.LockContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	session.Unlock()
}

type sessionRoundTripper func(*http.Request) (*http.Response, error)

func (f sessionRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestV2ClientCanceledWaiterDoesNotExecute(t *testing.T) {
	holderEntered := make(chan struct{})
	releaseHolder := make(chan struct{})
	var calls atomic.Int32
	transport := sessionRoundTripper(func(req *http.Request) (*http.Response, error) {
		call := calls.Add(1)
		if req.URL.Path == "/holder" && call == 1 {
			close(holderEntered)
			<-releaseHolder
			return &http.Response{StatusCode: http.StatusServiceUnavailable, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("retry"))}, nil
		}
		if req.Header.Get(HeaderCLISessionId) != "synthetic-session" {
			t.Error("missing session header")
		}
		return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})
	serverURL, _ := url.Parse("https://example.invalid")
	client := NewV2ClientWithHttpClient(&http.Client{Transport: transport}, serverURL, &RetryConfig{Enabled: true, RetryMax: 1, RetryWaitMin: time.Millisecond, RetryWaitMax: time.Millisecond})
	client.session = &Session{SessionId: "synthetic-session"}
	holderDone := make(chan error, 1)
	go func() {
		resp, err := client.doGetRequest(context.Background(), "/holder")
		if resp != nil {
			_ = resp.Body.Close()
		}
		holderDone <- err
	}()
	<-holderEntered
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(releaseHolder) }) }
	defer release()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	waiterDone := make(chan error, 1)
	go func() {
		resp, err := client.doGetRequest(ctx, "/canceled")
		if resp != nil {
			_ = resp.Body.Close()
		}
		waiterDone <- err
	}()
	select {
	case err := <-waiterDone:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected deadline, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("request waited for the holder after its deadline")
	}
	if calls.Load() != 1 {
		t.Fatal("canceled waiter reached transport")
	}
	release()
	if err := <-holderDone; err != nil {
		t.Fatal(err)
	}
	resp, err := client.doGetRequest(context.Background(), "/valid")
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if calls.Load() != 3 {
		t.Fatalf("expected holder, its retry, and valid request only, got %d", calls.Load())
	}
}

func TestSessionLockContextSerializesConcurrentCallers(t *testing.T) {
	session := &Session{}
	var active atomic.Int32
	var wg sync.WaitGroup
	for range 32 {
		wg.Go(func() {
			for range 20 {
				if err := session.LockContext(context.Background()); err != nil {
					t.Error(err)
					return
				}
				if active.Add(1) != 1 {
					t.Error("overlapping session holders")
				}
				active.Add(-1)
				session.Unlock()
			}
		})
	}
	wg.Wait()
}
