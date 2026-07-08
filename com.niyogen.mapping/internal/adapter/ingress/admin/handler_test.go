package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.niyogen/openclaw/internal/adapter/repository/memory"
	"com.niyogen/openclaw/internal/domain"
)

func TestRegisterApp(t *testing.T) {
	appRepo := memory.NewAppRepository()
	tkRepo := memory.NewTriggerKeyRepo()
	h := NewHandler(appRepo, tkRepo)

	reqBody := RegisterAppRequest{
		AppID:       "app-123",
		Name:        "Test App",
		Description: "A test application",
		BaseURL:     "http://example.com",
		AuthType:    "bearer",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/apps", bytes.NewReader(bodyBytes))
	req.Header.Set("X-Niyogen-Tenant", "tenant1")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Result().StatusCode)
	}

	app, err := appRepo.GetApp(context.Background(), "tenant1", "app-123")
	if err != nil {
		t.Fatalf("expected app to be stored, got error: %v", err)
	}
	if !app.Enabled {
		t.Errorf("expected app to be enabled by default")
	}
}

func TestManageTriggerKeys(t *testing.T) {
	appRepo := memory.NewAppRepository()
	tkRepo := memory.NewTriggerKeyRepo()
	h := NewHandler(appRepo, tkRepo)

	// Pre-register an app
	appRepo.RegisterApp(context.Background(), domain.RegisteredApp{
		AppID:    "app-123",
		TenantID: "tenant1",
	})

	reqBody := TriggerKeyRequest{
		IntentName:  "test_intent",
		UserID:      "user1",
		Scope:       "read_only",
		CallbackURL: "http://example.com/callback",
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/apps/app-123/trigger-keys", bytes.NewReader(bodyBytes))
	req.Header.Set("X-Niyogen-Tenant", "tenant1")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d. Body: %s", w.Result().StatusCode, w.Body.String())
	}

	var resp TriggerKeyResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.TriggerKey == "" {
		t.Errorf("expected a trigger key in response")
	}
}

func TestKillSwitch(t *testing.T) {
	appRepo := memory.NewAppRepository()
	tkRepo := memory.NewTriggerKeyRepo()
	h := NewHandler(appRepo, tkRepo)

	// Pre-register an app
	appRepo.RegisterApp(context.Background(), domain.RegisteredApp{
		AppID:    "app-123",
		TenantID: "tenant1",
		Enabled:  true,
	})

	req := httptest.NewRequest(http.MethodDelete, "/v1/apps/app-123/enable", nil)
	req.Header.Set("X-Niyogen-Tenant", "tenant1")

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Result().StatusCode)
	}

	app, _ := appRepo.GetApp(context.Background(), "tenant1", "app-123")
	if app.Enabled {
		t.Errorf("expected app to be disabled")
	}
}
