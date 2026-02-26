package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"usermanagement/internal/domain/user"
	"usermanagement/internal/infrastructure/logger"
)

// UserRepository implements domain.UserRepository using PostgreSQL.
type UserRepository struct {
	pool   *pgxpool.Pool
	logger *logger.Logger
}

// NewUserRepository creates a new PostgreSQL user repository.
func NewUserRepository(pool *pgxpool.Pool, logger *logger.Logger) *UserRepository {
	return &UserRepository{
		pool:   pool,
		logger: logger,
	}
}

// Save persists a new user.
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, query,
		u.ID(),
		u.Name(),
		u.Email(),
		u.CreatedAt(),
		u.UpdatedAt(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return user.ErrEmailExists
		}
		r.logger.Error("failed to save user", zap.Error(err))
		return fmt.Errorf("%w: %v", user.ErrRepositoryInternal, err)
	}

	return nil
}

// FindByID retrieves a user by ID.
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var uid uuid.UUID
	var name, email string
	var createdAt, updatedAt time.Time

	err := row.Scan(&uid, &name, &email, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		r.logger.Error("failed to find user by id", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", user.ErrRepositoryInternal, err)
	}

	return user.Reconstruct(uid, name, email, createdAt, updatedAt), nil
}

// FindByEmail retrieves a user by email.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.pool.QueryRow(ctx, query, email)

	var uid uuid.UUID
	var name, dbEmail string
	var createdAt, updatedAt time.Time

	err := row.Scan(&uid, &name, &dbEmail, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		r.logger.Error("failed to find user by email", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", user.ErrRepositoryInternal, err)
	}

	return user.Reconstruct(uid, name, dbEmail, createdAt, updatedAt), nil
}

// FindAll retrieves paginated users.
func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]*user.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("failed to list users", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", user.ErrRepositoryInternal, err)
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var uid uuid.UUID
		var name, email string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&uid, &name, &email, &createdAt, &updatedAt); err != nil {
			r.logger.Error("failed to scan user row", zap.Error(err))
			return nil, fmt.Errorf("%w: %v", user.ErrRepositoryInternal, err)
		}

		users = append(users, user.Reconstruct(uid, name, email, createdAt, updatedAt))
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error iterating user rows", zap.Error(err))
		return nil, fmt.Errorf("%w: %v", user.ErrRepositoryInternal, err)
	}

	return users, nil
}

// Update modifies an existing user.
func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.pool.Exec(ctx, query,
		u.Name(),
		u.Email(),
		u.UpdatedAt(),
		u.ID(),
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return user.ErrEmailExists
		}
		r.logger.Error("failed to update user", zap.Error(err))
		return fmt.Errorf("%w: %v", user.ErrRepositoryInternal, err)
	}

	if result.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

// Delete removes a user by ID.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error("failed to delete user", zap.Error(err))
		return fmt.Errorf("%w: %v", user.ErrRepositoryInternal, err)
	}

	if result.RowsAffected() == 0 {
		return user.ErrUserNotFound
	}

	return nil
}