package port

import (
	"context"

	"com.niyogen/openclaw/internal/domain"
)

// AppRepository manages external application registrations.
type AppRepository interface {
	RegisterApp(ctx context.Context, app domain.RegisteredApp) error
	GetApp(ctx context.Context, tenantID, appID string) (*domain.RegisteredApp, error)
	ListApps(ctx context.Context, tenantID string) ([]domain.RegisteredApp, error)
	UpdateApp(ctx context.Context, app domain.RegisteredApp) error
	DeleteApp(ctx context.Context, tenantID, appID string) error
}

// AuditRepository handles the immutable audit trail (Layer 7).
type AuditRepository interface {
	LogExecution(ctx context.Context, record AuditRecord) error
	ListRecords(ctx context.Context, tenantID string) ([]AuditRecord, error)
}

type AuditRecord struct {
	TenantID        string
	UserID          string
	Channel         string
	TrustTier       int
	Intent          string
	ParamsHash      string
	AppID           string
	Status          string // "success", "rejected", "error"
	RejectionReason string
	LatencyMs       int
	ResponseSummary string
}

// ReplayStore manages nonce tracking for replay protection.
type ReplayStore interface {
	// CheckAndStore rejects duplicate nonces within the TTL window.
	// The nonce is bound to channelType+senderID — a nonce from WhatsApp
	// user A cannot be replayed as an API call from user B.
	CheckAndStore(ctx context.Context, channelType, senderID, nonce string, ttlSeconds int) error
}
