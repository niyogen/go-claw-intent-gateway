package memory

import (
	"context"
	"fmt"
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
			err := store.CheckAndStore(ctx, "api", fmt.Sprintf("sender-%d", id), "nonce-1", 10)
			if err != nil {
				t.Errorf("goroutine %d: unexpected error: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
}

func TestReplayStore_TTL(t *testing.T) {
	store := NewReplayStore()
	ctx := context.Background()

	err := store.CheckAndStore(ctx, "whatsapp", "user-1", "nonce-1", 1)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Immediate reuse should fail
	err = store.CheckAndStore(ctx, "whatsapp", "user-1", "nonce-1", 1)
	if err != ErrNonceReused {
		t.Errorf("expected ErrNonceReused, got %v", err)
	}

	// Wait for TTL to expire
	time.Sleep(1500 * time.Millisecond)

	// After expiration, reuse should succeed
	err = store.CheckAndStore(ctx, "whatsapp", "user-1", "nonce-1", 1)
	if err != nil {
		t.Errorf("expected nil after expiration, got %v", err)
	}
}

func TestReplayStore_CrossChannelReplay(t *testing.T) {
	store := NewReplayStore()
	ctx := context.Background()

	err := store.CheckAndStore(ctx, "whatsapp", "user-1", "nonce-1", 10)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	// Same nonce, different channel -> should succeed
	err = store.CheckAndStore(ctx, "api", "user-1", "nonce-1", 10)
	if err != nil {
		t.Errorf("expected nil for different channel, got %v", err)
	}

	// Same nonce, same channel, different user -> should succeed
	err = store.CheckAndStore(ctx, "whatsapp", "user-2", "nonce-1", 10)
	if err != nil {
		t.Errorf("expected nil for different user, got %v", err)
	}
}

func TestValidateTimestamp(t *testing.T) {
	now := time.Now()
	maxAge := 5 * time.Minute
	maxFuture := 30 * time.Second

	// Fresh
	err := ValidateTimestamp(now.Add(-1*time.Minute), maxAge, maxFuture)
	if err != nil {
		t.Errorf("expected nil for fresh timestamp, got %v", err)
	}

	// Too old
	err = ValidateTimestamp(now.Add(-6*time.Minute), maxAge, maxFuture)
	if err == nil {
		t.Errorf("expected error for old timestamp, got nil")
	}

	// Future
	err = ValidateTimestamp(now.Add(1*time.Minute), maxAge, maxFuture)
	if err == nil {
		t.Errorf("expected error for future timestamp, got nil")
	}
}
