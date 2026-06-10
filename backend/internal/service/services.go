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
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Auth:          authService,
		Job:           s.Job,
		UploadService: NewUploadService(s, *repos.Upload),
	}, nil
}
