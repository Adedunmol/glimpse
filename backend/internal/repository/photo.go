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

func (r *PhotoRepository) CreatePhotos(ctx context.Context, userID string, payload *photo.CompletePhotosPayload) error {
	storageKeys := make([]string, len(payload.Files))
	for i, f := range payload.Files {
		storageKeys[i] = f.Key
	}

	stmt := `
		INSERT INTO photos (upload_id, storage_key, status)
		SELECT u.id, key, @status
 		FROM uploads u
 		CROSS JOIN unnest(@storage_keys::text[]) AS key
 		WHERE u.id = @upload_id AND u.host_id = @host_id
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"upload_id":    payload.UploadID,
		"storage_keys": storageKeys,
		"status":       "uploaded",
		"host_id":      userID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute batch create photos query for upload_id=%s user_id=%s: %w", payload.UploadID, userID, err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("upload not found or not owned by user (upload_id=%s user_id=%s)", payload.UploadID, userID)
	}

	return nil
}
