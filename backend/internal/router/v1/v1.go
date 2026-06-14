package v1

import (
	"github.com/Adedunmol/glimpse/internal/handler"
	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(router *echo.Group, handlers *handler.Handlers, middleware *middleware.Middlewares) {

	registerUploadRoutes(router, handlers.Upload, middleware.Auth)
	registerClusterRoutes(router, handlers.Cluster, middleware.Auth)
	registerLinkRoutes(router, handlers.Link, middleware.Auth)
}
