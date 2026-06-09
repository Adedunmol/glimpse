package cluster_photo

import (
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/google/uuid"
)

type ClusterPhoto struct {
	model.BaseWithCreatedAt
	model.BaseWithUpdatedAt

	ClusterID uuid.UUID `json:"clusterId" db:"cluster_id"`
	PhotoID   uuid.UUID `json:"photoId" db:"photo_id"`
}
