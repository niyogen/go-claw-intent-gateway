package memory

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestReplayStore_ConcurrentAccess(t *testing.T) {
	store := NewReplayStore()
	ctx := context.Background()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Should succeed first time
			_ = store.CheckAndStore(ctx, string(rune('A'+id)), 10)
		}(i)
	}

	wg.Wait()
}

func TestReplayStore_TTL(t *testing.T) {
	store := NewReplayStore()
	ctx := context.Background()

	err := store.CheckAndStore(ctx, "nonce-1", 1)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Immediate reuse should fail
	err = store.CheckAndStore(ctx, "nonce-1", 1)
	if err != ErrNonceReused {
		t.Errorf("expected ErrNonceReused, got %v", err)
	}

	// Wait for TTL to expire
	time.Sleep(1500 * time.Millisecond)

	// After expiration, reuse should succeed
	err = store.CheckAndStore(ctx, "nonce-1", 1)
	if err != nil {
		t.Errorf("expected nil after expiration, got %v", err)
	}
}
