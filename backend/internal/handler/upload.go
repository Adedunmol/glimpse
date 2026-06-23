package handler

import (
	"net/http"

	"github.com/Adedunmol/glimpse/internal/model/photo"

	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/upload"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/Adedunmol/glimpse/internal/service"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	Handler
	uploadService *service.UploadService
}

func NewUploadHandler(s *server.Server, uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		Handler:       NewHandler(s),
		uploadService: uploadService,
	}
}

func (h *UploadHandler) CreateUpload(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *upload.CreateUploadPayload) (*upload.Upload, error) {
			userID := middleware.GetUserID(c)
			return h.uploadService.CreateUpload(c, userID, payload)
		},
		http.StatusCreated,
		&upload.CreateUploadPayload{},
	)(c)
}

func (h *UploadHandler) GetUploadByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *upload.GetUploadByIDPayload) (*upload.PopulatedUpload, error) {
			userID := middleware.GetUserID(c)
			return h.uploadService.GetUploadByID(c, userID, payload.ID)
		},
		http.StatusOK,
		&upload.GetUploadByIDPayload{},
	)(c)
}

func (h *UploadHandler) GetUploads(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *upload.GetUploadsQuery) (*model.PaginatedResponse[upload.Upload], error) {
			userID := middleware.GetUserID(c)
			return h.uploadService.GetUploads(c, userID, query)
		},
		http.StatusOK,
		&upload.GetUploadsQuery{},
	)(c)
}

func (h *UploadHandler) UpdateUpload(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *upload.UpdateUploadPayload) (*upload.Upload, error) {
			userID := middleware.GetUserID(c)
			return h.uploadService.UpdateUpload(c, userID, payload)
		},
		http.StatusOK,
		&upload.UpdateUploadPayload{},
	)(c)
}

func (h *UploadHandler) DeleteUpload(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *upload.DeleteUploadPayload) error {
			userID := middleware.GetUserID(c)
			return h.uploadService.DeleteUpload(c, userID, payload.ID)
		},
		http.StatusNoContent,
		&upload.DeleteUploadPayload{},
	)(c)
}

func (h *UploadHandler) GetPresignedURLs(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *photo.CreatePhotosPayload) (*photo.PresignedURL, error) {
			userID := middleware.GetUserID(c)
			return h.uploadService.GetPresignedUrls(c, userID, payload)
		},
		http.StatusOK,
		&photo.CreatePhotosPayload{},
	)(c)
}

func (h *UploadHandler) CompleteUpload(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *photo.CompletePhotosPayload) error {
			userID := middleware.GetUserID(c)
			return h.uploadService.CompletePhotosUpload(c, userID, payload)
		},
		http.StatusNoContent,
		&photo.CompletePhotosPayload{},
	)(c)
}
