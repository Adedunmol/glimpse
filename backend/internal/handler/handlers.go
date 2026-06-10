package handler

import (
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/Adedunmol/glimpse/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
	Upload  *UploadHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
		Upload:  NewUploadHandler(s, services.UploadService),
	}
}
