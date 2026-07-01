package memory

import (
	"context"
	"sync"
	"testing"

	"com.niyogen/openclaw/internal/port"
)

func TestAuditRepo_ConcurrentAccess(t *testing.T) {
	repo := NewAuditRepository()
	ctx := context.Background()

	var wg sync.WaitGroup
	// Concurrently log records
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repo.LogExecution(ctx, port.AuditRecord{
				TenantID: "tenant-1",
			})
		}()
	}

	wg.Wait()

	// Type assertion to verify length since port doesn't expose it
	memRepo, ok := repo.(*auditRepo)
	if !ok {
		t.Fatalf("expected *auditRepo")
	}

	memRepo.mu.Lock()
	defer memRepo.mu.Unlock()
	if len(memRepo.records) != 100 {
		t.Errorf("expected 100 records, got %d", len(memRepo.records))
	}
}
