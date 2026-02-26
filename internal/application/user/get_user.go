package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"usermanagement/internal/domain/user"
)

// GetUserUseCase implements the get user use case.
type GetUserUseCase struct {
	repo user.UserRepository
}

// NewGetUserUseCase creates a new instance.
func NewGetUserUseCase(repo user.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{repo: repo}
}

// Execute retrieves a user by ID.
func (uc *GetUserUseCase) Execute(ctx context.Context, id uuid.UUID) (*UserOutput, error) {
	domainUser, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	output := MapFromDomain(domainUser)
	return &output, nil
}