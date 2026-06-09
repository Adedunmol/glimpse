package upload

import (
	"time"

	"github.com/Adedunmol/glimpse/internal/model"
)

type UploadStatus string

const (
	UploadStatusPending    UploadStatus = "pending"
	UploadStatusProcessing UploadStatus = "processing"
	UploadStatusDone       UploadStatus = "done"
	UploadStatusFailed     UploadStatus = "failed"
)

type Upload struct {
	model.Base

	UserID    string       `json:"userId" db:"user_id"`
	HostID    string       `json:"hostId" db:"host_id"`
	Status    UploadStatus `json:"status" db:"status"`
	ExpiresAt time.Time    `json:"expiresAt,omitempty" db:"expires_at"`
}
