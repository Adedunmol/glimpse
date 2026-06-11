package service

import (
	"fmt"
	"github.com/Adedunmol/glimpse/internal/lib/aws"
	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/photo"
	"github.com/Adedunmol/glimpse/internal/model/upload"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"path/filepath"
)

type UploadService struct {
	server     *server.Server
	uploadRepo repository.UploadRepository
	awsClient  *aws.AWS
}

func NewUploadService(server *server.Server, uploadRepo repository.UploadRepository, awsClient *aws.AWS) *UploadService {
	return &UploadService{
		server:     server,
		uploadRepo: uploadRepo,
		awsClient:  awsClient,
	}
}

func (s *UploadService) CreateUpload(ctx echo.Context, userID string, payload *upload.CreateUploadPayload) (*upload.Upload, error) {
	logger := middleware.GetLogger(ctx)

	uploadItem, err := s.uploadRepo.CreateUpload(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create upload")
		return nil, err
	}

	// business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "upload_created").
		Str("upload_id", "").
		Str("name", "").
		Msg("Upload created successfully")

	return uploadItem, nil
}

func (s *UploadService) GetUploadByID(ctx echo.Context, userID string, uploadID uuid.UUID) (*upload.Upload, error) {
	logger := middleware.GetLogger(ctx)

	uploadItem, err := s.uploadRepo.GetUploadByID(ctx.Request().Context(), userID, uploadID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch upload by ID")
		return nil, err
	}

	return uploadItem, nil
}

func (s *UploadService) GetUploads(ctx echo.Context, userID string, query *upload.GetUploadsQuery) (*model.PaginatedResponse[upload.Upload], error) {
	logger := middleware.GetLogger(ctx)

	uploads, err := s.uploadRepo.GetUploads(ctx.Request().Context(), userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch uploads")
		return nil, err
	}

	return uploads, nil
}

func (s *UploadService) UpdateUpload(ctx echo.Context, userID string, payload *upload.UpdateUploadPayload) (*upload.Upload, error) {
	logger := middleware.GetLogger(ctx)

	upload, err := s.uploadRepo.UpdateUpload(ctx.Request().Context(), userID, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update upload")
		return nil, err
	}

	// business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "upload_updated").
		Str("upload_id", upload.ID.String()).
		Str("name", upload.Name).
		Msg("Upload updated successfully")

	return upload, nil
}

func (s *UploadService) DeleteUpload(ctx echo.Context, userID string, uploadID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	err := s.uploadRepo.DeleteUpload(ctx.Request().Context(), userID, uploadID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete upload")
		return err
	}

	// business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "upload_deleted").
		Str("upload_id", uploadID.String()).
		Msg("Upload deleted successfully")

	return nil
}

func (s *UploadService) GetPresignedUrls(ctx echo.Context, userID string, payload *photo.CreatePhotosPayload) (*photo.PresignedURL, error) {
	logger := middleware.GetLogger(ctx)

	uploads := make([]photo.Upload, 0, len(payload.Files))

	for _, file := range payload.Files {
		key := fmt.Sprintf("users/%s/photos/%s%s", userID, uuid.New().String(), filepath.Ext(file.Name))

		url, err := s.awsClient.S3.CreatePresignedUploadURL(ctx.Request().Context(), s.server.Config.AWS.UploadBucket, key)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create presigned url")
			return nil, err
		}
		uploads = append(uploads, photo.Upload{
			Key: key,
			Url: url,
		})
	}

	result := &photo.PresignedURL{
		UploadID: payload.UploadID,
		Uploads:  uploads,
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "presigned_urls_generation").
		Str("upload_id", payload.UploadID).
		Msg("Presigned URLs generated successfully")

	return result, nil
}
