package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"com.niyogen/openclaw/internal/domain"
	"com.niyogen/openclaw/internal/port"
)

var ErrAppNotFound = errors.New("app not found")

type appRepo struct {
	mu   sync.RWMutex
	apps map[string]map[string]domain.RegisteredApp
}

func NewAppRepository() port.AppRepository {
	return &appRepo{
		apps: make(map[string]map[string]domain.RegisteredApp),
	}
}

func (r *appRepo) RegisterApp(ctx context.Context, app domain.RegisteredApp) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.apps[app.TenantID]; !ok {
		r.apps[app.TenantID] = make(map[string]domain.RegisteredApp)
	}
	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	if !app.Enabled {
		app.Enabled = true // default to enabled on register
	}

	r.apps[app.TenantID][app.AppID] = app
	return nil
}

func (r *appRepo) GetApp(ctx context.Context, tenantID, appID string) (*domain.RegisteredApp, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenantApps, ok := r.apps[tenantID]
	if !ok {
		return nil, ErrAppNotFound
	}
	app, ok := tenantApps[appID]
	if !ok {
		return nil, ErrAppNotFound
	}
	return &app, nil
}

func (r *appRepo) ListApps(ctx context.Context, tenantID string) ([]domain.RegisteredApp, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tenantApps, ok := r.apps[tenantID]
	if !ok {
		return []domain.RegisteredApp{}, nil
	}

	apps := make([]domain.RegisteredApp, 0, len(tenantApps))
	for _, app := range tenantApps {
		apps = append(apps, app)
	}
	return apps, nil
}

func (r *appRepo) UpdateApp(ctx context.Context, app domain.RegisteredApp) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantApps, ok := r.apps[app.TenantID]
	if !ok {
		return ErrAppNotFound
	}
	if _, ok := tenantApps[app.AppID]; !ok {
		return ErrAppNotFound
	}
	tenantApps[app.AppID] = app
	return nil
}

func (r *appRepo) DeleteApp(ctx context.Context, tenantID, appID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantApps, ok := r.apps[tenantID]
	if !ok {
		return ErrAppNotFound
	}
	if _, ok := tenantApps[appID]; !ok {
		return ErrAppNotFound
	}
	delete(tenantApps, appID)
	return nil
}

// UpdateAppEnabled updates the kill switch state for an app.
func (r *appRepo) UpdateAppEnabled(ctx context.Context, tenantID, appID string, enabled bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tenantApps, ok := r.apps[tenantID]
	if !ok {
		return fmt.Errorf("app not found")
	}

	app, ok := tenantApps[appID]
	if !ok {
		return fmt.Errorf("app not found")
	}

	app.Enabled = enabled
	app.UpdatedAt = time.Now()
	tenantApps[appID] = app

	return nil
}
