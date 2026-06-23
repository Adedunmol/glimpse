package cluster

import (
	"time"

	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Cluster struct {
	model.Base

	UploadID         string    `json:"uploadId" db:"upload_id"`
	Label            string    `json:"label" db:"label"`
	ThumbnailPhotoId uuid.UUID `json:"thumbnailPhotoId" db:"thumbnail_photo_id"`
}

type CreateLinkCommand struct {
	ClusterID uuid.UUID  `param:"clusterId" validate:"required,uuid"`
	Password  string     `json:"password,omitempty" validate:"omitempty,min=6"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

func (p *CreateLinkCommand) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
