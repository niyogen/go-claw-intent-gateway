package domain

import "time"

// RegisteredApp represents an external application registered with OpenClaw.
type RegisteredApp struct {
	AppID       string
	TenantID    string
	Name        string
	Description string
	BaseURL     string
	AuthType    string
	ConfigHash       string // SHA256 of the YAML config for tamper detection
	Endpoints        []IntentDef
	Events           []EventDef
	AllowedSourceIPs []string
	RequireTLS       bool
	Active           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// EventDef represents an event that the app can emit
type EventDef struct {
	Name          string
	Description   string
	MaxRate       string // e.g. "30/minute"
	PayloadSchema map[string]ParamRule
}

// IntentScope defines if an intent mutates state
type IntentScope string

const (
	ScopeReadOnly IntentScope = "read_only"
	ScopeWrite    IntentScope = "write"
)

// IntentDef represents a single capability an app exposes.
type IntentDef struct {
	Name                 string
	Method               string
	Path                 string
	Description          string
	Permission           Role
	ReadOnly             bool
	RequiresConfirmation bool
	ParamsSchema         map[string]ParamRule
	TimeoutSeconds       int
}

// ParamRule defines the schema constraint for an intent parameter.
type ParamRule struct {
	Type      string `json:"type"`     // string, int, bool
	Required  bool   `json:"required"`
	MaxLength int    `json:"max_length,omitempty"`
}

// TriggerKey represents a signed token to asynchronously trigger an intent.
type TriggerKey struct {
	KeyHash     string
	IntentName  string
	UserID      string
	TenantID    string
	Scope       IntentScope
	CallbackURL string
	ExpiresAt   time.Time
}
