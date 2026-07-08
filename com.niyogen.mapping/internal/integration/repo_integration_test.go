package integration_test

import (
	"context"
	"testing"
	"time"

	"com.niyogen/openclaw/internal/adapter/repository/memory"
	"com.niyogen/openclaw/internal/domain"
	"com.niyogen/openclaw/internal/port"
)

func TestRepositoryIntegration(t *testing.T) {
	ctx := context.Background()

	appRepo := memory.NewAppRepository()
	auditRepo := memory.NewAuditRepository()
	replayStore := memory.NewReplayStore()

	// 1. App Registration
	app := domain.RegisteredApp{
		AppID:       "app-int-1",
		TenantID:    "tenant-1",
		Name:        "Integration App",
		Description: "An integration test application",
		BaseURL:     "https://api.int.testapp.com",
		AuthType:    "Bearer",
		ConfigHash:  "hash_int",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := appRepo.RegisterApp(ctx, app)
	if err != nil {
		t.Fatalf("failed to register app: %v", err)
	}

	// 2. Fetch App
	fetchedApp, err := appRepo.GetApp(ctx, "tenant-1", "app-int-1")
	if err != nil {
		t.Fatalf("failed to get app: %v", err)
	}
	if fetchedApp.Name != "Integration App" {
		t.Errorf("expected Integration App, got %s", fetchedApp.Name)
	}

	// 3. Replay Protection
	nonceKey := "nonce-int-123"
	err = replayStore.CheckAndStore(ctx, "whatsapp", "user-1", nonceKey, 60)
	if err != nil {
		t.Fatalf("first replay check failed: %v", err)
	}

	// Second check should fail
	err = replayStore.CheckAndStore(ctx, "whatsapp", "user-1", nonceKey, 60)
	if err == nil {
		t.Errorf("expected error on replay, got nil")
	}

	// 4. Audit Logging
	auditRecord := port.AuditRecord{
		TenantID:        "tenant-1",
		UserID:          "user-1",
		Channel:         "whatsapp",
		TrustTier:       3,
		Intent:          "create_user",
		ParamsHash:      "param_hash",
		AppID:           "app-int-1",
		Status:          "success",
		RejectionReason: "",
		LatencyMs:       45,
		ResponseSummary: "Created successfully",
	}

	err = auditRepo.LogExecution(ctx, auditRecord)
	if err != nil {
		t.Fatalf("failed to log audit execution: %v", err)
	}

	records, _ := auditRepo.ListRecords(ctx, "tenant-1")
	if len(records) != 1 {
		t.Errorf("expected 1 record, got %d", len(records))
	}

	// 5. App Enable/Disable (Kill Switch)
	err = appRepo.UpdateAppEnabled(ctx, "tenant-1", "app-int-1", false)
	if err != nil {
		t.Fatalf("failed to update app enabled state: %v", err)
	}
	disabledApp, _ := appRepo.GetApp(ctx, "tenant-1", "app-int-1")
	if disabledApp.Enabled {
		t.Errorf("expected app to be disabled")
	}

	// 6. Trigger Keys
	tkRepo := memory.NewTriggerKeyRepo()
	tk := domain.TriggerKey{
		KeyHash:     "hash123",
		IntentName:  "test_intent",
		UserID:      "user1",
		TenantID:    "tenant-1",
		Scope:       domain.ScopeReadOnly,
		CallbackURL: "http://example.com/callback",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}
	err = tkRepo.SaveTriggerKey(ctx, tk)
	if err != nil {
		t.Fatalf("failed to save trigger key: %v", err)
	}
	fetchedTk, err := tkRepo.GetTriggerKey(ctx, "hash123")
	if err != nil || fetchedTk.IntentName != "test_intent" {
		t.Fatalf("failed to fetch trigger key or mismatch intent: %v", err)
	}
}
