package repository

import (
	"context"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/model/photo"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/jackc/pgx/v5"
)

type PhotoRepository struct {
	server *server.Server
}

func NewPhotoRepository(server *server.Server) *PhotoRepository {
	return &PhotoRepository{
		server: server,
	}
}

func (r *PhotoRepository) CreatePhotos(ctx context.Context, userID string, payload *photo.CompletePhotosPayload) ([]photo.Photo, error) {
	storageKeys := make([]string, len(payload.Files))
	for i, f := range payload.Files {
		storageKeys[i] = f.Key
	}

	stmt := `
		INSERT INTO photos (upload_id, storage_key, status)
		SELECT @upload_id, key, @status
		FROM unnest(@storage_keys::text[]) AS key
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"upload_id":    payload.UploadID,
		"storage_keys": storageKeys,
		"status":       "uploaded",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute batch create photos query for upload_id=%s user_id=%s: %w", payload.UploadID, userID, err)
	}

	photos, err := pgx.CollectRows(rows, pgx.RowToStructByName[photo.Photo])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows from table:photos for upload_id=%s user_id=%s: %w", payload.UploadID, userID, err)
	}

	return photos, nil
}
