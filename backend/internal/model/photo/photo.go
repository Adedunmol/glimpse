package photo

import "github.com/Adedunmol/glimpse/internal/model"

type Status string

const (
	PhotoStatusPending  Status = "pending"
	PhotoStatusUploaded Status = "uploaded"
)

type Photo struct {
	model.Base
	UploadID   string `json:"uploadId" db:"upload_id"`
	StorageKey string `json:"storageKey" db:"storage_key"`
	Status     Status `json:"status" db:"status"`
}
