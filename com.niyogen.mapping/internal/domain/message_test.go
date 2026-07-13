package domain_test

import (
	"testing"
	"time"

	"com.niyogen/openclaw/internal/domain"
)

func TestUnifiedMessage_Initialization(t *testing.T) {
	msg := domain.UnifiedMessage{
		ChannelType:     domain.ChannelWhatsApp,
		ChannelID:       "+1234567890",
		SenderID:        "whatsapp-user-id",
		TenantID:        "acme-corp",
		UserID:          "john-doe",
		Role:            domain.RoleViewer,
		TrustTier:       domain.Tier3,
		Text:            "Hello, world!",
		Nonce:           "unique-nonce-123",
		Timestamp:       time.Now(),
		ChannelSigValid: true,
	}

	if msg.ChannelType != domain.ChannelWhatsApp {
		t.Errorf("expected ChannelWhatsApp, got %s", msg.ChannelType)
	}
	if msg.TrustTier != domain.Tier3 {
		t.Errorf("expected Tier3, got %d", msg.TrustTier)
	}
}
