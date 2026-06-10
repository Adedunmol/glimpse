package repository

import (
	"context"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/model/user"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/jackc/pgx/v5"
)

var ErrNotFound = pgx.ErrNoRows

type UserRepository interface {
	GetUserEmail(ctx context.Context, email string) (string, error)
	GetUserByID(ctx context.Context, userID string) (*user.User, error)
	CreateUser(ctx context.Context, email, clerkID string) (*user.User, error)
	DeleteUser(ctx context.Context, userId string) error
}

type PostgresUserRepository struct {
	server *server.Server
}

func NewPostgresRepository(s *server.Server) *PostgresUserRepository {
	return &PostgresUserRepository{
		server: s,
	}
}

func (p *PostgresUserRepository) GetUserEmail(ctx context.Context, email string) (string, error) {
	stmt := `
		SELECT
			*
		FROM
			users
		WHERE
			email=@email
	`
	rows, err := p.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"email": email,
	})

	if err != nil {
		return "", fmt.Errorf("failed to execute query to fetch user email: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return "", fmt.Errorf("failed to collect row from table:users: %w", err)
	}

	return user.Email, nil
}

func (p *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*user.User, error) {
	stmt := `
			SELECT
				*
			FROM
				users
			WHERE
				user_id=@userID
	`
	rows, err := p.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"email": userID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute query to fetch user email: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users: %w", err)
	}

	return &user, nil
}

func (p *PostgresUserRepository) CreateUser(ctx context.Context, email, clerkId string) (*user.User, error) {
	stmt := `
		INSERT INTO
			users (user_id, email, created_at, updated_at)
		VALUES
			(
				@user_id, @email
			)
		RETURNING
			*
	`
	rows, err := p.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": clerkId,
		"email":   email,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute create user query: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, fmt.Errorf("failed too collect row from table:user: %w", err)
	}

	return &user, nil
}

func (p *PostgresUserRepository) DeleteUser(ctx context.Context, userId string) error {
	result, err := p.server.DB.Pool.Exec(ctx, `
		DELETE FROM users WHERE user_id = @id
	`, pgx.NamedArgs{
		"id": userId,
	})

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
