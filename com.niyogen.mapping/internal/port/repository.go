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
}

// AuditRepository handles the immutable audit trail (Layer 7).
type AuditRepository interface {
	LogExecution(ctx context.Context, record AuditRecord) error
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
	// CheckAndStore returns an error if the nonce was already seen.
	CheckAndStore(ctx context.Context, key string, ttlSeconds int) error
}
