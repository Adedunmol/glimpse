package photo

import (
	"github.com/Adedunmol/glimpse/internal/model"
)

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
	IsEmbedded bool   `json:"isEmbedded" db:"is_embedded"`
}

type PresignedURL struct {
	UploadID string   `json:"uploadId"`
	Uploads  []Upload `json:"uploads"`
}

type Upload struct {
	Key string `json:"key"`
	Url string `json:"url"`
}
