package user

import (
	"context"
	"fmt"

	"usermanagement/internal/domain/user"
)

// UpdateUserUseCase implements the update user use case.
type UpdateUserUseCase struct {
	repo user.UserRepository
}

// NewUpdateUserUseCase creates a new instance.
func NewUpdateUserUseCase(repo user.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

// Execute updates a user.
func (uc *UpdateUserUseCase) Execute(ctx context.Context, input UpdateUserInput) (*UserOutput, error) {
	// Retrieve existing
	domainUser, err := uc.repo.FindByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	// Check email uniqueness if changing email
	if input.Email != nil && *input.Email != domainUser.Email() {
		existing, err := uc.repo.FindByEmail(ctx, *input.Email)
		if err != nil && !errors.Is(err, user.ErrUserNotFound) {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if existing != nil {
			return nil, user.ErrEmailExists
		}
		
		if err := domainUser.UpdateEmail(*input.Email); err != nil {
			return nil, err
		}
	}

	// Update name if provided
	if input.Name != nil {
		if err := domainUser.UpdateName(*input.Name); err != nil {
			return nil, err
		}
	}

	// Persist
	if err := uc.repo.Update(ctx, domainUser); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	output := MapFromDomain(domainUser)
	return &output, nil
}