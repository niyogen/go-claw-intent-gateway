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

	records, err := repo.ListRecords(ctx, "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 100 {
		t.Errorf("expected 100 records, got %d", len(records))
	}
}
