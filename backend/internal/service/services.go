package service

import (
	"github.com/Adedunmol/go-boilerplate/internal/lib/job"
	"github.com/Adedunmol/go-boilerplate/internal/repository"
	"github.com/Adedunmol/go-boilerplate/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Auth: authService,
		Job:  s.Job,
	}, nil
}
