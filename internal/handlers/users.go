package handlers

import (
	"context"
	"net/http"

	"monolith/internal/database"
	"monolith/internal/models"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	db *database.DB
}

func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	var users []models.User
	err := pgxscan.Select(context.Background(), h.db.Pool, &users, `
		SELECT id, username, email, name, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at 
		FROM account 
		WHERE is_disabled = FALSE
		ORDER BY created_at DESC
	`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	err := pgxscan.Get(context.Background(), h.db.Pool, &user, `
		SELECT id, username, email, name, avatar, is_admin, language, theme, timezone, 
		       last_seen_at, is_disabled, created_at, updated_at 
		FROM account 
		WHERE id = $1 AND is_disabled = FALSE
	`, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"required,email"`
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	var user models.User
	err := pgxscan.Get(context.Background(), h.db.Pool, &user, `
		INSERT INTO account (username, name, email, created_at, updated_at) 
		VALUES ($1, $2, $3, NOW(), NOW()) 
		RETURNING id, username, email, name, is_admin, language, theme, timezone, 
		          last_seen_at, is_disabled, created_at, updated_at
	`, req.Username, req.Name, req.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	var user models.User
	err := pgxscan.Get(context.Background(), h.db.Pool, &user, `
		UPDATE account 
		SET username = $1, name = $2, email = $3, updated_at = NOW() 
		WHERE id = $4 AND is_disabled = FALSE
		RETURNING id, username, email, name, is_admin, language, theme, timezone, 
		          last_seen_at, is_disabled, created_at, updated_at
	`, req.Username, req.Name, req.Email, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	_, err := h.db.Pool.Exec(context.Background(), `
		UPDATE account SET is_disabled = TRUE, updated_at = NOW() WHERE id = $1
	`, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
