package service

import (
	"github.com/Adedunmol/glimpse/internal/lib/job"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
)

type Services struct {
	Auth          *AuthService
	Job           *job.JobService
	UploadService *UploadService
	ClerkService  *ClerkService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	clerkService := NewClerkService(s, repos.UserRepository)

	return &Services{
		Auth:          authService,
		Job:           s.Job,
		UploadService: NewUploadService(s, *repos.Upload),
		ClerkService:  clerkService,
	}, nil
}
