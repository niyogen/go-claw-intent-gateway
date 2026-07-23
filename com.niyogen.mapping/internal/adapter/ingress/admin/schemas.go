package admin

import "com.niyogen/openclaw/internal/domain"

// RegisterAppRequest represents the request to register an application.
// @Summary Register a new application
// @Description Register an external app to expose its intents to OpenClaw.
// @Accept json
// @Produce json
// @Param request body RegisterAppRequest true "Application details"
type RegisterAppRequest struct {
	AppID       string `json:"app_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BaseURL     string `json:"base_url"`
	AuthType    string `json:"auth_type"`
}

// RegisterAppResponse represents the response upon successful registration.
// @Summary Application registration response
// @Produce json
type RegisterAppResponse struct {
	Status string               `json:"status"`
	App    domain.RegisteredApp `json:"app"`
}

// TriggerKeyRequest represents a request to generate a trigger key.
// @Summary Generate a Trigger Key
// @Description Create a new trigger key for an application to initiate async flows.
// @Accept json
// @Produce json
type TriggerKeyRequest struct {
	IntentName  string `json:"intent_name"`
	UserID      string `json:"user_id"`
	Scope       string `json:"scope"`
	CallbackURL string `json:"callback_url"`
}

// TriggerKeyResponse is the response containing the generated trigger key.
// @Summary Trigger Key response
// @Produce json
type TriggerKeyResponse struct {
	TriggerKey string `json:"trigger_key"`
	ExpiresAt  string `json:"expires_at"`
}

// SuccessResponse is a generic response for successful operations.
// @Produce json
type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// ErrorResponse is a generic error response.
// @Produce json
type ErrorResponse struct {
	Error string `json:"error"`
}
