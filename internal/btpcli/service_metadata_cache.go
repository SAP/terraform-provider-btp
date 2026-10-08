package btpcli

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// DefaultServiceMetadataCacheTTL limits catalogue staleness to five minutes.
const DefaultServiceMetadataCacheTTL = 5 * time.Minute
const serviceMetadataCacheCapacity = 1024

type serviceMetadataKey struct {
	session *Session
	request string
}
type serviceMetadataEntry struct {
	data        []byte
	status      int
	contentType string
	expires     time.Time
}
type serviceMetadataCache struct {
	mu      sync.Mutex
	ttl     time.Duration
	now     func() time.Time
	entries map[serviceMetadataKey]serviceMetadataEntry
	pending map[serviceMetadataKey]chan struct{}
}

func newServiceMetadataCache(ttl time.Duration) *serviceMetadataCache {
	return &serviceMetadataCache{ttl: ttl, now: time.Now, entries: make(map[serviceMetadataKey]serviceMetadataEntry), pending: make(map[serviceMetadataKey]chan struct{})}
}

// SetServiceMetadataCacheTTL configures catalogue caching before using the client.
// Zero disables caching. Changing the TTL clears all previous entries.
func (v2 *v2Client) SetServiceMetadataCacheTTL(ttl time.Duration) {
	v2.serviceMetadataCache = newServiceMetadataCache(ttl)
}

// doExecuteServiceMetadata caches only successful individual plan/offering reads.
// The request includes the subaccount, lookup kind and every argument; the cache
// belongs to one client, and a replacement login session cannot reuse old entries.
// JSON snapshots ensure callers never share mutable catalogue maps or slices.
func doExecuteServiceMetadata[T any](client *v2Client, ctx context.Context, req *CommandRequest) (T, CommandResponse, error) {
	var obj T
	if err := ctx.Err(); err != nil {
		return obj, CommandResponse{}, err
	}
	cache := client.serviceMetadataCache
	if cache == nil || cache.ttl <= 0 {
		return doExecute[T](client, ctx, req)
	}
	encoded, err := json.Marshal(req)
	if err != nil {
		return obj, CommandResponse{}, err
	}
	key := serviceMetadataKey{session: client.session, request: string(encoded)}
	for {
		cache.mu.Lock()
		now := cache.now()
		if entry, ok := cache.entries[key]; ok {
			if now.Before(entry.expires) {
				cache.mu.Unlock()
				err := json.Unmarshal(entry.data, &obj)
				return obj, CommandResponse{StatusCode: entry.status, ContentType: entry.contentType}, err
			}
			delete(cache.entries, key)
		}
		if pending, ok := cache.pending[key]; ok {
			cache.mu.Unlock()
			select {
			case <-ctx.Done():
				return obj, CommandResponse{}, ctx.Err()
			case <-pending:
			}
			if err := ctx.Err(); err != nil {
				return obj, CommandResponse{}, err
			}
			continue
		}
		// Bound concurrent bookkeeping too. Excess distinct lookups bypass caching.
		if len(cache.pending) >= serviceMetadataCacheCapacity {
			cache.mu.Unlock()
			return doExecute[T](client, ctx, req)
		}
		pending := make(chan struct{})
		cache.pending[key] = pending
		cache.mu.Unlock()
		result, response, requestErr := doExecute[T](client, ctx, req)
		data, marshalErr := json.Marshal(result)
		cache.mu.Lock()
		if requestErr == nil && marshalErr == nil && ctx.Err() == nil {
			now = cache.now()
			for k, entry := range cache.entries {
				if !now.Before(entry.expires) {
					delete(cache.entries, k)
				}
			}
			if len(cache.entries) >= serviceMetadataCacheCapacity {
				var oldest serviceMetadataKey
				var expires time.Time
				for k, entry := range cache.entries {
					if expires.IsZero() || entry.expires.Before(expires) {
						oldest, expires = k, entry.expires
					}
				}
				delete(cache.entries, oldest)
			}
			cache.entries[key] = serviceMetadataEntry{data: data, status: response.StatusCode, contentType: response.ContentType, expires: now.Add(cache.ttl)}
		}
		delete(cache.pending, key)
		close(pending)
		cache.mu.Unlock()
		return result, response, requestErr
	}
}
