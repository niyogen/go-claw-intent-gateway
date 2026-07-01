package memory

import (
	"context"
	"sync"

	"com.niyogen/openclaw/internal/port"
)

type auditRepo struct {
	mu      sync.Mutex
	records []port.AuditRecord
}

func NewAuditRepository() port.AuditRepository {
	return &auditRepo{
		records: make([]port.AuditRecord, 0),
	}
}

func (r *auditRepo) LogExecution(ctx context.Context, record port.AuditRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = append(r.records, record)
	return nil
}
