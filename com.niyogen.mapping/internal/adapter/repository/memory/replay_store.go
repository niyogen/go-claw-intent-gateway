package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"com.niyogen/openclaw/internal/port"
)

var ErrNonceReused = errors.New("nonce already used")

type replayStore struct {
	mu     sync.Mutex
	nonces map[string]time.Time
}

func NewReplayStore() port.ReplayStore {
	return &replayStore{
		nonces: make(map[string]time.Time),
	}
}

func (r *replayStore) CheckAndStore(ctx context.Context, key string, ttlSeconds int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Lazy cleanup: clean expired keys (optional but good for memory)
	for k, expiry := range r.nonces {
		if now.After(expiry) {
			delete(r.nonces, k)
		}
	}

	expiry, exists := r.nonces[key]
	if exists {
		if now.Before(expiry) {
			return ErrNonceReused
		}
		// If it exists but is expired, it would have been deleted above, 
		// but in case it's the current key, we can reuse it, although the loop above 
		// should already have cleared it.
	}

	r.nonces[key] = now.Add(time.Duration(ttlSeconds) * time.Second)
	return nil
}
