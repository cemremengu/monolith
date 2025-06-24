package api

import (
	"net/http"

	"monolith/internal/database"
	"monolith/internal/service/user"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *user.Service
}

func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{
		userService: user.NewService(db),
	}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.userService.GetUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")

	user, err := h.userService.GetUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req user.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	createdUser, err := h.userService.CreateUser(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, createdUser)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	var req user.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	updatedUser, err := h.userService.UpdateUser(c.Request().Context(), id, req)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")

	err := h.userService.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
