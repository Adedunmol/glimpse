package repository

import (
	"github.com/Adedunmol/glimpse/internal/lib/notification"
	"github.com/Adedunmol/glimpse/internal/server"
)

type Repositories struct {
	Upload         *UploadRepository
	Photo          *PhotoRepository
	UserRepository UserRepository
	Cluster        *ClusterRepository
	Link           *LinkRepository
	Device         *notification.DeviceRepository
}

func NewRepositories(s *server.Server, deviceRepo *notification.DeviceRepository) *Repositories {
	userRepo := NewPostgresRepository(s)
	return &Repositories{
		UserRepository: userRepo,
		Upload:         NewUploadRepository(s),
		Photo:          NewPhotoRepository(s),
		Cluster:        NewClusterRepository(s),
		Link:           NewLinkRepository(s),
		Device:         deviceRepo,
	}
}
