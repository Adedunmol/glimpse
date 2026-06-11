package service

import (
	"fmt"
	"github.com/Adedunmol/glimpse/internal/lib/aws"
	"github.com/Adedunmol/glimpse/internal/lib/job"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
)

type Services struct {
	Auth         *AuthService
	Job          *job.JobService
	ClerkService *ClerkService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)
	clerkService := NewClerkService(s, repos.UserRepository)
	_, err := aws.NewAWS(s) // TODO: embed in photo service
	if err != nil {
		return nil, fmt.Errorf("error creating aws client: %w", err)
	}

	return &Services{
		Auth:         authService,
		Job:          s.Job,
		ClerkService: clerkService,
	}, nil
}
