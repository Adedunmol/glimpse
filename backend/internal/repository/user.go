package repository

import (
	"context"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetUserEmail(ctx context.Context, email string) (string, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	CreateUser(ctx context.Context, email, clerkID string) (*model.User, error)
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
		return "", fmt.Errorf("failed to execute query to fetch user email")
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return "", fmt.Errorf("failed to collect row from table:users")
	}

	return user.Email, nil
}

func (p *PostgresUserRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
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
		return nil, fmt.Errorf("failed to execute query to fetch user email")
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:users")
	}

	return &user, nil
}

func (p *PostgresUserRepository) CreateUser(ctx context.Context, email, clerkId string) (*model.User, error) {
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
		return nil, fmt.Errorf("failed to execute create user query")
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nil, fmt.Errorf("failed too collect row from table:user")
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
