package cluster

import (
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/google/uuid"
)

type Cluster struct {
	model.Base

	UploadID         string    `json:"uploadId" db:"upload_id"`
	Label            string    `json:"label" db:"label"`
	ThumbnailPhotoId uuid.UUID `json:"thumbnailPhotoId" db:"thumbnail_photo_id"`
}
