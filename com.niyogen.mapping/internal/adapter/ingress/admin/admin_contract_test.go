package admin_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"com.niyogen/openclaw/internal/adapter/ingress/admin"
	"com.niyogen/openclaw/internal/adapter/repository/memory"
)

func TestAdminAPIContract(t *testing.T) {
	appRepo := memory.NewAppRepository()
	tkRepo := memory.NewTriggerKeyRepo()
	handler := admin.NewHandler(appRepo, tkRepo)
	server := httptest.NewServer(handler)
	defer server.Close()

	t.Run("POST /v1/apps - Register App", func(t *testing.T) {
		reqBody := map[string]string{
			"app_id":      "test-app-1",
			"name":        "Test App",
			"description": "A test application",
			"base_url":    "https://api.test.com",
			"auth_type":   "hmac",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/v1/apps", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Niyogen-Tenant", "tenant-1")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		var respData map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&respData)
		if respData["status"] != "success" {
			t.Errorf("Expected status success, got %v", respData["status"])
		}
	})

	t.Run("POST /v1/apps/{id}/trigger-keys - Generate Trigger Key", func(t *testing.T) {
		reqBody := map[string]string{
			"intent_name":  "create_ticket",
			"user_id":      "user-1",
			"scope":        "read_write",
			"callback_url": "https://callback.com",
		}
		bodyBytes, _ := json.Marshal(reqBody)
		
		req, _ := http.NewRequest(http.MethodPost, server.URL+"/v1/apps/test-app-1/trigger-keys", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Niyogen-Tenant", "tenant-1")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /v1/apps/{id}/enable - Kill Switch", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, server.URL+"/v1/apps/test-app-1/enable", nil)
		req.Header.Set("X-Niyogen-Tenant", "tenant-1")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})
}
