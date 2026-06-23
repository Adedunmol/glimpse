package service

import (
	"fmt"

	"github.com/Adedunmol/glimpse/internal/lib/aws"
	"github.com/Adedunmol/glimpse/internal/lib/job"
	"github.com/Adedunmol/glimpse/internal/lib/notification"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
)

type Services struct {
	Auth           *AuthService
	UploadService  *UploadService
	ClerkService   *ClerkService
	ClusterService *ClusterService
	LinkService    *LinkService
	DeviceService  *notification.DeviceService
	JobService     *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories, deviceService *notification.DeviceService, jobService *job.JobService) (*Services, error) {
	authService := NewAuthService(s)
	clerkService := NewClerkService(s, repos.UserRepository)
	linkService := NewLinkService(s, repos.Link)
	//cluster service depends on link service to create links
	clusterService := NewClusterService(s, repos.Cluster, linkService)
	awsClient, err := aws.NewAWS(s) // TODO: embed in photo service
	if err != nil {
		return nil, fmt.Errorf("error creating aws client: %w", err)
	}
	// deviceService := NewDeviceService(s, repos.Device)
	return &Services{
		Auth:           authService,
		UploadService:  NewUploadService(s, *repos.Upload, *repos.Photo, awsClient),
		ClerkService:   clerkService,
		ClusterService: clusterService,
		LinkService:    linkService,
		DeviceService:  deviceService,
		JobService:     jobService,
	}, nil
}
