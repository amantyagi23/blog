package user

import (
	"time"
	"usermanagement/internal/domain/user"

	"github.com/google/uuid"
)

// DTOs decouple domain from delivery layer.
// They represent data contracts for the application layer.

// CreateUserInput represents data needed to create a user.
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserInput represents data needed to update a user.
type UpdateUserInput struct {
	ID    uuid.UUID `json:"-"` // From URL param, not body
	Name  *string   `json:"name,omitempty"`
	Email *string   `json:"email,omitempty"`
}

// UserOutput represents user data returned to clients.
type UserOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MapFromDomain converts domain entity to output DTO.
func MapFromDomain(u * user.User) UserOutput {
	return UserOutput{
		ID:        u.ID(),
		Name:      u.Name(),
		Email:     u.Email(),
		CreatedAt: u.CreatedAt(),
		UpdatedAt: u.UpdatedAt(),
	}
}

// PaginationInput for list operations.
type PaginationInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ListUsersOutput represents paginated user list.
type ListUsersOutput struct {
	Users []*UserOutput `json:"users"`
	Total int           `json:"total"`
}