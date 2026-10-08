package btpcli

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServiceMetadataCacheFreshnessAndScope(t *testing.T) {
	var calls atomic.Int32
	client, server := prepareClientFacadeForTest(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		_, _ = fmt.Fprintf(w, `{"id":"plan-%d","name":"original"}`, n)
	})
	defer server.Close()
	now := time.Unix(1000, 0)
	client.serviceMetadataCache.now = func() time.Time { return now }
	ctx := context.Background()
	first, _, err := client.Services.Plan.GetById(ctx, "account-a", "plan-a")
	require.NoError(t, err)
	again, _, err := client.Services.Plan.GetById(ctx, "account-a", "plan-a")
	require.NoError(t, err)
	require.Equal(t, first, again)
	require.EqualValues(t, 1, calls.Load())
	now = now.Add(DefaultServiceMetadataCacheTTL)
	_, _, err = client.Services.Plan.GetById(ctx, "account-a", "plan-a")
	require.NoError(t, err)
	require.EqualValues(t, 2, calls.Load())
	_, _, err = client.Services.Plan.GetById(ctx, "account-b", "plan-a")
	require.NoError(t, err)
	_, _, err = client.Services.Plan.GetByName(ctx, "account-a", "plan-a", "offering-a")
	require.NoError(t, err)
	_, _, err = client.Services.Offering.GetById(ctx, "account-a", "plan-a")
	require.NoError(t, err)
	require.EqualValues(t, 5, calls.Load())
	client.session = &Session{}
	_, _, err = client.Services.Plan.GetById(ctx, "account-a", "plan-a")
	require.NoError(t, err)
	require.EqualValues(t, 6, calls.Load())
	client.SetServiceMetadataCacheTTL(0)
	for range 2 {
		_, _, err = client.Services.Plan.GetById(ctx, "account-a", "plan-a")
		require.NoError(t, err)
	}
	require.EqualValues(t, 8, calls.Load())
}

func TestServiceMetadataCacheCopiesAndDoesNotCacheErrors(t *testing.T) {
	var calls atomic.Int32
	client, server := prepareClientFacadeForTest(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			w.Header().Set(HeaderCLIBackendStatus, "503")
			_, _ = fmt.Fprint(w, `{"error":"temporary"}`)
			return
		}
		_, _ = fmt.Fprint(w, `{"id":"offering-a","metadata":{"displayName":"blue"}}`)
	})
	defer server.Close()
	ctx := context.Background()
	_, _, err := client.Services.Offering.GetById(ctx, "account-a", "offering-a")
	require.Error(t, err)
	first, _, err := client.Services.Offering.GetById(ctx, "account-a", "offering-a")
	require.NoError(t, err)
	first.Metadata.DisplayName = "red"
	again, _, err := client.Services.Offering.GetById(ctx, "account-a", "offering-a")
	require.NoError(t, err)
	require.Equal(t, "blue", again.Metadata.DisplayName)
	require.EqualValues(t, 2, calls.Load())
	ctx, cancel := context.WithCancel(ctx)
	cancel()
	_, _, err = client.Services.Offering.GetById(ctx, "account-a", "offering-a")
	require.ErrorIs(t, err, context.Canceled)
}

func TestServiceMetadataCacheCollapsesConcurrentMissesAndCancelsWaiters(t *testing.T) {
	var calls atomic.Int32
	entered, release := make(chan struct{}), make(chan struct{})
	client, server := prepareClientFacadeForTest(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			close(entered)
		}
		<-release
		_, _ = fmt.Fprint(w, `{"id":"plan-a"}`)
	})
	defer server.Close()
	var wg sync.WaitGroup
	errors := make(chan error, 10)
	wg.Go(func() {
		_, _, err := client.Services.Plan.GetById(context.Background(), "account-a", "plan-a")
		errors <- err
	})
	<-entered
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, _, err := client.Services.Plan.GetById(ctx, "account-a", "plan-a")
	require.ErrorIs(t, err, context.DeadlineExceeded)
	for range 9 {
		wg.Go(func() {
			_, _, err := client.Services.Plan.GetById(context.Background(), "account-a", "plan-a")
			errors <- err
		})
	}
	close(release)
	wg.Wait()
	close(errors)
	for err := range errors {
		require.NoError(t, err)
	}
	require.EqualValues(t, 1, calls.Load())
}

func TestServiceMetadataCacheBounded(t *testing.T) {
	client, server := prepareClientFacadeForTest(func(w http.ResponseWriter, r *http.Request) { _, _ = fmt.Fprint(w, `{"id":"plan-a"}`) })
	defer server.Close()
	for i := range serviceMetadataCacheCapacity + 5 {
		_, _, err := client.Services.Plan.GetById(context.Background(), "account-a", fmt.Sprint(i))
		require.NoError(t, err)
	}
	require.Len(t, client.serviceMetadataCache.entries, serviceMetadataCacheCapacity)
}
