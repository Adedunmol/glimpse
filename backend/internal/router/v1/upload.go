package v1

import (
	"github.com/Adedunmol/glimpse/internal/handler"
	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerUploadRoutes(r *echo.Group, h *handler.UploadHandler, auth *middleware.AuthMiddleware) {

	uploads := r.Group("/uploads")
	uploads.Use(auth.RequireAuth)

	uploads.POST("", h.CreateUpload)
	uploads.GET("", h.GetUploads)

	dynamicUpload := uploads.Group("/:id")
	dynamicUpload.GET("", h.GetUploadByID)
	dynamicUpload.PATCH("", h.UpdateUpload)
	dynamicUpload.DELETE("", h.DeleteUpload)

	photos := dynamicUpload.Group("/photos")
	photos.POST("", h.GetPresignedURLs)
}
