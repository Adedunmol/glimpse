package service

import (
	"fmt"

	"github.com/Adedunmol/glimpse/internal/lib/aws"
	"github.com/Adedunmol/glimpse/internal/lib/fcm"
	"github.com/Adedunmol/glimpse/internal/lib/job"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
)

type Services struct {
	Auth                *AuthService
	Job                 *job.JobService
	UploadService       *UploadService
	ClerkService        *ClerkService
	ClusterService      *ClusterService
	LinkService         *LinkService
	DeviceService       *DeviceService
	NotificationService *NotificationService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	clerkService := NewClerkService(s, repos.UserRepository)
	clusterService := NewClusterService(s, repos.Cluster)
	linkService := NewLinkService(s, repos.Link)
	awsClient, err := aws.NewAWS(s) // TODO: embed in photo service
	if err != nil {
		return nil, fmt.Errorf("error creating aws client: %w", err)
	}
	deviceService := NewDeviceService(s, repos.Device)
	credPath := s.Config.FCM.CredentialPath
	projectID := s.Config.FCM.ProjectID
	fcmClient, err := fcm.NewFCMClient(credPath, projectID)
	if err != nil {
		return nil, fmt.Errorf("error creating fcm client: %w", err)
	}
	notificationService := NewNotificationService(s, repos.Device, fcmClient)

	return &Services{
		Auth:                authService,
		Job:                 s.Job,
		UploadService:       NewUploadService(s, *repos.Upload, *repos.Photo, awsClient),
		ClerkService:        clerkService,
		ClusterService:      clusterService,
		LinkService:         linkService,
		DeviceService:       deviceService,
		NotificationService: notificationService,
	}, nil
}
