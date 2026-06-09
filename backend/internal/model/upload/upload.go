package upload

import (
	"time"

	"github.com/Adedunmol/glimpse/internal/model"
)

type Status string

const (
	UploadStatusPending    Status = "pending"
	UploadStatusProcessing Status = "processing"
	UploadStatusDone       Status = "done"
	UploadStatusFailed     Status = "failed"
)

type Upload struct {
	model.Base

	Name      string    `json:"name" db:"name"`
	HostID    string    `json:"hostId" db:"host_id"`
	Status    Status    `json:"status" db:"status"`
	ExpiresAt time.Time `json:"expiresAt,omitempty" db:"expires_at"`
}
