package memory

import (
	"context"
	"fmt"
	"sync"

	"com.niyogen/openclaw/internal/domain"
	"com.niyogen/openclaw/internal/port"
)

type triggerKeyRepo struct {
	mu   sync.RWMutex
	keys map[string]domain.TriggerKey // keyHash -> TriggerKey
}

func NewTriggerKeyRepo() port.TriggerKeyRepository {
	return &triggerKeyRepo{
		keys: make(map[string]domain.TriggerKey),
	}
}

func (r *triggerKeyRepo) SaveTriggerKey(ctx context.Context, key domain.TriggerKey) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.keys[key.KeyHash] = key
	return nil
}

func (r *triggerKeyRepo) GetTriggerKey(ctx context.Context, keyHash string) (*domain.TriggerKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	key, ok := r.keys[keyHash]
	if !ok {
		return nil, fmt.Errorf("trigger key not found")
	}
	return &key, nil
}

func (r *triggerKeyRepo) RevokeTriggerKey(ctx context.Context, keyHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.keys[keyHash]; !ok {
		return fmt.Errorf("trigger key not found")
	}
	delete(r.keys, keyHash)
	return nil
}

func (r *triggerKeyRepo) ListTriggerKeys(ctx context.Context, tenantID string) ([]domain.TriggerKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domain.TriggerKey
	for _, key := range r.keys {
		if key.TenantID == tenantID {
			result = append(result, key)
		}
	}
	return result, nil
}
