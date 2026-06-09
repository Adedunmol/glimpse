package link

import (
	"time"

	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/google/uuid"
)

type Link struct {
	model.Base

	ClusterID           uuid.UUID `json:"clusterId" db:"cluster_id"`
	Token               string    `json:"token" db:"token"`
	IsPasswordProtected bool      `json:"isPasswordProtected" db:"is_password_protected"`
	PasswordHash        *string   `json:"passwordHash" db:"password_hash"`
	ExpiresAt           time.Time `json:"expiresAt" db:"expires_at"`
	IsActive            bool      `json:"isActive" db:"is_active"`
}
