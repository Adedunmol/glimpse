package photo

import "github.com/Adedunmol/glimpse/internal/model"

type PhotoStatus string

const (
	PhotoStatusPending  PhotoStatus = "pending"
	PhotoStatusUploaded PhotoStatus = "uploaded"
)

type Photo struct {
	model.Base
	UploadID   string      `json:"uploadId" db:"upload_id"`
	StorageKey string      `json:"storageKey" db:"storage_key"`
	Status     PhotoStatus `json:"status" db:"status"`
}
