package repository

import (
	"context"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/model/user_device"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/jackc/pgx/v5"
)

type DeviceRepository struct {
	server *server.Server
}

func NewDeviceRepository(server *server.Server) *DeviceRepository {
	return &DeviceRepository{
		server: server,
	}
}

func (r *DeviceRepository) UpsertDevice(ctx context.Context, userID, fcmToken, platform string) (*user_device.UserDevice, error) {
	stmt := `
		INSERT INTO
			user_devices (user_id, push_token, platform)
		VALUES
			(@userId, @pushToken, @platform)
		ON CONFLICT
			(push_token)
		DO UPDATE SET
			user_id = @userId, platform = @platform
		RETURNING
			*
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"userId":    userID,
		"pushToken": fcmToken,
		"platform":  platform,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute create devuce query for user_id=%s", userID)
	}

	deviceItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user_device.UserDevice])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:device_token for user_id=%s", userID)
	}
	return &deviceItem, nil
}

func (r *DeviceRepository) GetUserTokens(ctx context.Context, userID string) ([]user_device.UserDevice, error) {
	stmt := `
		SELECT
			*
		FROM
			user_devices
		WHERE
			user_id = @userId
		LIMIT 100
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"userId": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute fetch user devices for user_id=%s", userID)
	}

	devices, err := pgx.CollectRows(rows, pgx.RowToStructByName[user_device.UserDevice])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows from table:user_devices for user_id=%s", userID)
	}
	return devices, nil
}
