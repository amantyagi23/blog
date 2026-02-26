package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"usermanagement/internal/domain/user"
)

// DeleteUserUseCase implements the delete user use case.
type DeleteUserUseCase struct {
	repo user.UserRepository
}

// NewDeleteUserUseCase creates a new instance.
func NewDeleteUserUseCase(repo user.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{repo: repo}
}

// Execute deletes a user.
func (uc *DeleteUserUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	// Verify existence first
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return user.ErrUserNotFound
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}