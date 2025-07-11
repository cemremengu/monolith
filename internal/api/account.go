package api

import (
	"errors"
	"net/http"

	"monolith/internal/service/account"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	accountService *account.Service
}

func NewAccountHandler(accountService *account.Service) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) Profile(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, APIError{Message: "Invalid user ID"})
	}

	account, err := h.accountService.GetAccountByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, APIError{Message: "Account not found"})
	}

	return c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) UpdatePreferences(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return c.JSON(http.StatusUnauthorized, APIError{Message: "Invalid user ID"})
	}

	var req account.UpdatePreferencesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Message: "Invalid request body"})
	}

	updatedAccount, err := h.accountService.UpdatePreferences(c.Request().Context(), userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIError{Message: "Failed to update preferences"})
	}

	return c.JSON(http.StatusOK, updatedAccount)
}

func (h *AccountHandler) Register(c echo.Context) error {
	var req account.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Message: "Invalid request body"})
	}

	userAccount, err := h.accountService.Register(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrPasswordTooShort):
			return c.JSON(http.StatusBadRequest, APIError{Message: err.Error()})
		case errors.Is(err, account.ErrUserAlreadyExists):
			return c.JSON(http.StatusConflict, APIError{Message: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, APIError{Message: "Failed to register account"})
		}
	}

	response := map[string]any{
		"account": userAccount,
	}

	return c.JSON(http.StatusCreated, response)
}
