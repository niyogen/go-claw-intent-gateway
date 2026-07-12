package domain_test

import (
	"testing"
	"time"

	"com.niyogen/openclaw/internal/domain"
)

func TestRegisteredApp_Initialization(t *testing.T) {
	now := time.Now()
	app := domain.RegisteredApp{
		AppID:       "app-1",
		TenantID:    "tenant-1",
		Name:        "Test App",
		Description: "A test application",
		BaseURL:     "https://api.testapp.com",
		AuthType:    "Bearer",
		ConfigHash:  "hash123",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if app.AppID != "app-1" {
		t.Errorf("expected app-1, got %s", app.AppID)
	}
	if app.TenantID != "tenant-1" {
		t.Errorf("expected tenant-1, got %s", app.TenantID)
	}
	if app.CreatedAt != now {
		t.Errorf("expected time %v, got %v", now, app.CreatedAt)
	}
	if app.Active != false {
		t.Errorf("expected Active to be false by default")
	}
}

func TestEventDef_Initialization(t *testing.T) {
	event := domain.EventDef{
		Name:        "inspection_failed",
		Description: "Triggered when inspection fails",
		MaxRate:     "30/minute",
		PayloadSchema: map[string]domain.ParamRule{
			"equipment_id": {Type: "string"},
		},
	}

	if event.Name != "inspection_failed" {
		t.Errorf("expected inspection_failed, got %s", event.Name)
	}
	if event.PayloadSchema["equipment_id"].Type != "string" {
		t.Errorf("expected string type for equipment_id")
	}
}

func TestIntentDef_Initialization(t *testing.T) {
	intent := domain.IntentDef{
		Name:                 "create_user",
		Method:               "POST",
		Path:                 "/users",
		Description:          "Creates a new user",
		Permission:           domain.RoleAdmin,
		ReadOnly:             false,
		RequiresConfirmation: true,
		ParamsSchema: map[string]domain.ParamRule{
			"username": {
				Type:      "string",
				Required:  true,
				MaxLength: 50,
			},
		},
		TimeoutSeconds: 30,
	}

	if intent.Name != "create_user" {
		t.Errorf("expected create_user, got %s", intent.Name)
	}
	if intent.Permission != domain.RoleAdmin {
		t.Errorf("expected admin, got %s", intent.Permission)
	}
	if intent.ParamsSchema["username"].Required != true {
		t.Errorf("expected true for required parameter")
	}
}

func TestTriggerKey_Initialization(t *testing.T) {
	now := time.Now()
	expires := now.Add(1 * time.Hour)
	key := domain.TriggerKey{
		KeyHash:     "key-hash-xyz",
		IntentName:  "update_status",
		UserID:      "user-1",
		TenantID:    "tenant-1",
		Scope:       domain.ScopeWrite,
		CallbackURL: "https://callback.com",
		ExpiresAt:   expires,
	}

	if key.Scope != domain.ScopeWrite {
		t.Errorf("expected ScopeWrite, got %s", key.Scope)
	}
	if key.ExpiresAt != expires {
		t.Errorf("expected time %v, got %v", expires, key.ExpiresAt)
	}
}
