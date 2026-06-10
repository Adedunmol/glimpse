package handler

import (
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/Adedunmol/glimpse/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
	Clerk   *ClerkWebHookHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	clerk, err := NewClerkWebHookHandler(s, services.ClerkService)
	if err != nil {
		s.Logger.Error().Err(err).Msg("error initializing clerk webhook handler")
		panic(err)
	}
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
		Clerk:   clerk,
	}
}
