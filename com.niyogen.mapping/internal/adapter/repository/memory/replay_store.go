package memory

import (
	"context"
	"errors"
	"fmt"
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

func (r *replayStore) CheckAndStore(ctx context.Context, channelType, senderID, nonce string, ttlSeconds int) error {
	key := channelType + ":" + senderID + ":" + nonce

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

func ValidateTimestamp(timestamp time.Time, maxAge, maxFuture time.Duration) error {
	age := time.Since(timestamp)
	if age > maxAge {
		return fmt.Errorf("message too old: %v exceeds max age %v", age, maxAge)
	}
	if age < -maxFuture {
		return fmt.Errorf("message from the future: %v exceeds max skew %v", -age, maxFuture)
	}
	return nil
}
