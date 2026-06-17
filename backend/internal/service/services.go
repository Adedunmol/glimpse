package service

import (
	"fmt"

	"github.com/Adedunmol/glimpse/internal/lib/aws"
	"github.com/Adedunmol/glimpse/internal/lib/job"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
)

type Services struct {
	Auth           *AuthService
	Job            *job.JobService
	UploadService  *UploadService
	ClerkService   *ClerkService
	ClusterService *ClusterService
	LinkService    *LinkService
	DeviceService  *DeviceService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	clerkService := NewClerkService(s, repos.UserRepository)
	clusterService := NewClusterService(s, repos.Cluster)
	linkService := NewLinkService(s, repos.Link)
	awsClient, err := aws.NewAWS(s) // TODO: embed in photo service
	deviceService := NewDeviceService(s, repos.Device)
	if err != nil {
		return nil, fmt.Errorf("error creating aws client: %w", err)
	}

	return &Services{
		Auth:           authService,
		Job:            s.Job,
		UploadService:  NewUploadService(s, *repos.Upload, *repos.Photo, awsClient),
		ClerkService:   clerkService,
		ClusterService: clusterService,
		LinkService:    linkService,
		DeviceService:  deviceService,
	}, nil
}
