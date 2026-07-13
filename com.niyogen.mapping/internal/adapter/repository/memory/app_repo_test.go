package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"com.niyogen/openclaw/internal/domain"
)

func TestAppRepo_ConcurrentAccess(t *testing.T) {
	repo := NewAppRepository()
	ctx := context.Background()

	var wg sync.WaitGroup
	// Concurrently register apps
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			repo.RegisterApp(ctx, domain.RegisteredApp{
				TenantID: "tenant-1",
				AppID:    fmt.Sprintf("app-%d", id),
			})
		}(i)
	}

	wg.Wait()

	apps, err := repo.ListApps(ctx, "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 100 {
		t.Errorf("expected 100 apps, got %d", len(apps))
	}
}

func TestAppRepo_GetApp(t *testing.T) {
	repo := NewAppRepository()
	ctx := context.Background()

	err := repo.RegisterApp(ctx, domain.RegisteredApp{
		TenantID: "tenant-1",
		AppID:    "app-1",
		Name:     "Test App",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	app, err := repo.GetApp(ctx, "tenant-1", "app-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.Name != "Test App" {
		t.Errorf("expected 'Test App', got '%s'", app.Name)
	}

	_, err = repo.GetApp(ctx, "tenant-1", "app-nonexistent")
	if err != ErrAppNotFound {
		t.Errorf("expected ErrAppNotFound, got %v", err)
	}
}

func TestAppRepo_UpdateApp(t *testing.T) {
	repo := NewAppRepository()
	ctx := context.Background()

	err := repo.RegisterApp(ctx, domain.RegisteredApp{
		TenantID: "tenant-1",
		AppID:    "app-1",
		Name:     "Test App",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = repo.UpdateApp(ctx, domain.RegisteredApp{
		TenantID: "tenant-1",
		AppID:    "app-1",
		Name:     "Updated App",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	app, _ := repo.GetApp(ctx, "tenant-1", "app-1")
	if app.Name != "Updated App" {
		t.Errorf("expected 'Updated App', got '%s'", app.Name)
	}

	// Update non-existent app
	err = repo.UpdateApp(ctx, domain.RegisteredApp{
		TenantID: "tenant-1",
		AppID:    "app-nonexistent",
	})
	if err != ErrAppNotFound {
		t.Errorf("expected ErrAppNotFound, got %v", err)
	}
}

func TestAppRepo_DeleteApp(t *testing.T) {
	repo := NewAppRepository()
	ctx := context.Background()

	repo.RegisterApp(ctx, domain.RegisteredApp{
		TenantID: "tenant-1",
		AppID:    "app-1",
	})

	err := repo.DeleteApp(ctx, "tenant-1", "app-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = repo.GetApp(ctx, "tenant-1", "app-1")
	if err != ErrAppNotFound {
		t.Errorf("expected ErrAppNotFound, got %v", err)
	}

	// Delete non-existent app
	err = repo.DeleteApp(ctx, "tenant-1", "app-nonexistent")
	if err != ErrAppNotFound {
		t.Errorf("expected ErrAppNotFound, got %v", err)
	}
}
