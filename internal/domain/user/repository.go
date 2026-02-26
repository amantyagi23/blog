package user

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository defines the contract for user persistence.
// It belongs to the domain layer - implementation details are in infrastructure.
// This is the OUTPUT PORT in Clean Architecture terminology.
type UserRepository interface {
	// Save persists a new user.
	Save(ctx context.Context, user *User) error
	
	// FindByID retrieves a user by their unique ID.
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	
	// FindByEmail retrieves a user by email (for uniqueness checks).
	FindByEmail(ctx context.Context, email string) (*User, error)
	
	// FindAll retrieves paginated users.
	FindAll(ctx context.Context, limit, offset int) ([]*User, error)
	
	// Update modifies an existing user.
	Update(ctx context.Context, user *User) error
	
	// Delete removes a user by ID.
	Delete(ctx context.Context, id uuid.UUID) error
}