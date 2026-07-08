package admin

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"com.niyogen/openclaw/internal/domain"
	"com.niyogen/openclaw/internal/port"
)

// Handler handles admin API endpoints for application registration and management.
type Handler struct {
	mux           *http.ServeMux
	appRepo       port.AppRepository
	triggerKeyRepo port.TriggerKeyRepository
}

// NewHandler creates a new admin HTTP handler.
func NewHandler(appRepo port.AppRepository, triggerKeyRepo port.TriggerKeyRepository) *Handler {
	h := &Handler{
		mux:            http.NewServeMux(),
		appRepo:        appRepo,
		triggerKeyRepo: triggerKeyRepo,
	}
	h.routes()
	return h
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) routes() {
	h.mux.HandleFunc("POST /v1/apps", h.handleRegisterApp)
	h.mux.HandleFunc("POST /v1/apps/{id}/trigger-keys", h.handleManageTriggerKeys)
	h.mux.HandleFunc("DELETE /v1/apps/{id}/enable", h.handleKillSwitch)
}

func (h *Handler) getTenantID(r *http.Request) string {
	tenant := r.Header.Get("X-Niyogen-Tenant")
	if tenant == "" {
		return "default-tenant" // fallback for simplicity in this phase
	}
	return tenant
}

func (h *Handler) handleRegisterApp(w http.ResponseWriter, r *http.Request) {
	var req RegisterAppRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	tenantID := h.getTenantID(r)

	app := domain.RegisteredApp{
		AppID:       req.AppID,
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		BaseURL:     req.BaseURL,
		AuthType:    req.AuthType,
		Enabled:     true,
	}

	if err := h.appRepo.RegisterApp(r.Context(), app); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to register app"})
		return
	}

	writeJSON(w, http.StatusCreated, RegisterAppResponse{Status: "success", App: app})
}

func (h *Handler) handleManageTriggerKeys(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("id")
	tenantID := h.getTenantID(r)

	// Verify app exists
	app, err := h.appRepo.GetApp(r.Context(), tenantID, appID)
	if err != nil || app == nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "app not found"})
		return
	}

	var req TriggerKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid request payload"})
		return
	}

	// Generate a random key hash (mocking the crypto step for simplicity here)
	bytes := make([]byte, 16)
	rand.Read(bytes)
	keyHash := hex.EncodeToString(bytes)

	tk := domain.TriggerKey{
		KeyHash:     keyHash,
		IntentName:  req.IntentName,
		UserID:      req.UserID,
		TenantID:    tenantID,
		Scope:       domain.IntentScope(req.Scope),
		CallbackURL: req.CallbackURL,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	if err := h.triggerKeyRepo.SaveTriggerKey(r.Context(), tk); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to generate trigger key"})
		return
	}

	writeJSON(w, http.StatusCreated, TriggerKeyResponse{
		TriggerKey: "tk_hmac_" + keyHash, // sending pseudo key
		ExpiresAt:  tk.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *Handler) handleKillSwitch(w http.ResponseWriter, r *http.Request) {
	appID := r.PathValue("id")
	tenantID := h.getTenantID(r)

	err := h.appRepo.UpdateAppEnabled(r.Context(), tenantID, appID, false)
	if err != nil {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "app not found or update failed"})
		return
	}

	writeJSON(w, http.StatusOK, SuccessResponse{Status: "success", Message: "app disabled"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
