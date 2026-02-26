package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	app "usermanagement/internal/application/user"
	"usermanagement/internal/domain/user"
	"usermanagement/internal/infrastructure/logger"
)

// UserHandler handles HTTP requests for user management.
type UserHandler struct {
	createUC *app.CreateUserUseCase
	getUC    *app.GetUserUseCase
	listUC   *app.ListUsersUseCase
	updateUC *app.UpdateUserUseCase
	deleteUC *app.DeleteUserUseCase
	logger   *logger.Logger
}

// NewUserHandler creates a new HTTP handler with injected use cases.
func NewUserHandler(
	createUC *app.CreateUserUseCase,
	getUC *app.GetUserUseCase,
	listUC *app.ListUsersUseCase,
	updateUC *app.UpdateUserUseCase,
	deleteUC *app.DeleteUserUseCase,
	logger *logger.Logger,
) *UserHandler {
	return &UserHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		logger:   logger,
	}
}

// Create handles POST /users.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input app.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	output, err := h.createUC.Execute(r.Context(), input)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, output)
}

// GetByID handles GET /users/{id}.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id format")
		return
	}

	output, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

// List handles GET /users with pagination.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parsePagination(r)

	input := app.PaginationInput{
		Limit:  limit,
		Offset: offset,
	}

	output, err := h.listUC.Execute(r.Context(), input)
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondJSON(w, http.StatusOK, output)
}

// Update handles PUT /users/{id}.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id format")
		return
	}

	var input app.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	input.ID = id

	output, err := h.updateUC.Execute(r.Context(), input)
	if err != nil {
		h.handleDomainError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, output)
}

// Delete handles DELETE /users/{id}.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id format")
		return
	}

	if err := h.deleteUC.Execute(r.Context(), id); err != nil {
		h.handleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleDomainError maps domain errors to HTTP status codes.
func (h *UserHandler) handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrUserNotFound):
		respondError(w, http.StatusNotFound, "user not found")
	case errors.Is(err, user.ErrEmailExists):
		respondError(w, http.StatusConflict, "email already exists")
	case errors.Is(err, user.ErrEmptyName):
		respondError(w, http.StatusBadRequest, "name cannot be empty")
	case errors.Is(err, user.ErrInvalidEmail):
		respondError(w, http.StatusBadRequest, "invalid email format")
	default:
		h.logger.Error("unexpected error", zap.Error(err))
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func parsePagination(r *http.Request) (limit, offset int) {
	query := r.URL.Query()
	
	limitStr := query.Get("limit")
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	} else {
		limit = 10
	}

	offsetStr := query.Get("offset")
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	} else {
		offset = 0
	}

	return
}