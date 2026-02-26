package user

import (
	"context"
	"fmt"

	"usermanagement/internal/domain/user"
)

// CreateUserUseCase implements the create user use case.
type CreateUserUseCase struct {
	repo user.UserRepository
}

// NewCreateUserUseCase creates a new instance.
func NewCreateUserUseCase(repo user.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

// Execute runs the use case.
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	// Check email uniqueness
	existing, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if existing != nil {
		return nil, user.ErrEmailExists
	}

	// Create domain entity (validates invariants)
	domainUser, err := user.New(input.Name, input.Email)
	if err != nil {
		return nil, err // Domain error propagates directly
	}

	// Persist
	if err := uc.repo.Save(ctx, domainUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	output := MapFromDomain(domainUser)
	return &output, nil
}