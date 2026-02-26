package user

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents the aggregate root of the User domain.
// It encapsulates business invariants and rules.
type User struct {
	id        uuid.UUID
	name      string
	email     string
	createdAt time.Time
	updatedAt time.Time
}

// Domain errors - part of the ubiquitous language
var (
	ErrEmptyName     = errors.New("user name cannot be empty")
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrNilUser       = errors.New("user cannot be nil")
	ErrUserNotFound  = errors.New("user not found")
	ErrEmailExists   = errors.New("email already exists")
)

// New creates a new User with validated invariants.
// This is the only way to create a valid User entity.
func New(name, email string) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrEmptyName
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		id:        uuid.New(),
		name:      strings.TrimSpace(name),
		email:     strings.ToLower(strings.TrimSpace(email)),
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Reconstruct rebuilds a User from persistence layer.
// Used by repositories when hydrating from database.
// Does NOT validate - assumes data is already valid from DB.
func Reconstruct(id uuid.UUID, name, email string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		name:      name,
		email:     email,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// UpdateName changes the user's name with validation.
func (u *User) UpdateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}
	u.name = strings.TrimSpace(name)
	u.updatedAt = time.Now().UTC()
	return nil
}

// UpdateEmail changes the user's email with validation.
func (u *User) UpdateEmail(email string) error {
	if err := validateEmail(email); err != nil {
		return err
	}
	u.email = strings.ToLower(strings.TrimSpace(email))
	u.updatedAt = time.Now().UTC()
	return nil
}

// ID returns the user's unique identifier.
func (u *User) ID() uuid.UUID {
	return u.id
}

// Name returns the user's name.
func (u *User) Name() string {
	return u.name
}

// Email returns the user's email.
func (u *User) Email() string {
	return u.email
}

// CreatedAt returns the creation timestamp.
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update timestamp.
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return ErrInvalidEmail
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}