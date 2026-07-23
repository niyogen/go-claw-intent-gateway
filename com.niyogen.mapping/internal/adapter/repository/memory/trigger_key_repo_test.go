package memory

import (
	"context"
	"testing"
	"time"

	"com.niyogen/openclaw/internal/domain"
)

func TestTriggerKeyRepo(t *testing.T) {
	repo := NewTriggerKeyRepo()
	ctx := context.Background()

	key := domain.TriggerKey{
		KeyHash:     "hash123",
		IntentName:  "test_intent",
		UserID:      "user1",
		TenantID:    "tenant1",
		Scope:       domain.ScopeReadOnly,
		CallbackURL: "http://example.com/callback",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	// Test Save
	err := repo.SaveTriggerKey(ctx, key)
	if err != nil {
		t.Fatalf("expected no error saving key, got %v", err)
	}

	// Test Get
	gotKey, err := repo.GetTriggerKey(ctx, "hash123")
	if err != nil {
		t.Fatalf("expected no error getting key, got %v", err)
	}
	if gotKey.IntentName != "test_intent" {
		t.Errorf("expected intent test_intent, got %s", gotKey.IntentName)
	}

	// Test List
	keys, err := repo.ListTriggerKeys(ctx, "tenant1")
	if err != nil {
		t.Fatalf("expected no error listing keys, got %v", err)
	}
	if len(keys) != 1 {
		t.Errorf("expected 1 key, got %d", len(keys))
	}

	// Test Revoke
	err = repo.RevokeTriggerKey(ctx, "hash123")
	if err != nil {
		t.Fatalf("expected no error revoking key, got %v", err)
	}

	_, err = repo.GetTriggerKey(ctx, "hash123")
	if err == nil {
		t.Error("expected error getting revoked key, got nil")
	}
}
