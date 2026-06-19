package domain

import "time"

// ChannelType enumerates the supported channels
type ChannelType string

const (
	ChannelWhatsApp ChannelType = "whatsapp"
	ChannelSlack    ChannelType = "slack"
	ChannelWebChat  ChannelType = "web_chat"
	ChannelAPI      ChannelType = "api"
	ChannelEmail    ChannelType = "email"
)

// TrustTier represents the security level of a channel
type TrustTier int

const (
	Tier1 TrustTier = 1 // Highest: Session cookie, mTLS
	Tier2 TrustTier = 2 // Medium: Slack signing secret
	Tier3 TrustTier = 3 // Lowest: WhatsApp phone mapping, Email
)

// Role represents the authorization level of a user
type Role string

const (
	RoleAdmin   Role = "admin"
	RoleEditor  Role = "editor"
	RoleViewer  Role = "viewer"
	RoleBlocked Role = "blocked"
)

// Attachment represents a file or media attached to a message
type Attachment struct {
	Type string // e.g. "image/png", "application/pdf"
	URL  string
	Size int64
}

// UnifiedMessage represents a normalized message from any channel
// as defined in RFC-001 (Layer 1-2).
type UnifiedMessage struct {
	// Identity (resolved by channel adapter)
	ChannelType ChannelType // "whatsapp" | "slack" | "web_chat" | "api" | "email"
	ChannelID   string      // channel-specific identifier
	SenderID    string      // user identity within the channel
	TenantID    string      // resolved from channel mappings
	UserID      string      // resolved from channel mappings
	Role        Role        // "admin" | "editor" | "viewer" | "blocked"
	TrustTier   TrustTier   // 1 (highest) to 3 (lowest)

	// Content
	Text        string       // the natural language message
	Attachments []Attachment // associated files

	// Security
	Nonce           string    // unique message ID for replay protection
	Timestamp       time.Time // message creation time
	ChannelSigValid bool      // true if channel-native signature verified

	// Context
	AppID          string // which app this is about (optional)
	ConversationID string // for multi-turn clarification
}
