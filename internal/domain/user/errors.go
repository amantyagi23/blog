package user

import "errors"

// Repository errors for infrastructure to use
var (
	ErrRepositoryConflict = errors.New("data conflict in repository")
	ErrRepositoryInternal = errors.New("internal repository error")
)