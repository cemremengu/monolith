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
		return APIError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid user ID",
		}
	}

	account, err := h.accountService.GetAccountByID(c.Request().Context(), userID)
	if err != nil {
		return APIError{
			Code:    http.StatusNotFound,
			Message: "Account not found",
			Err:     err,
		}
	}

	return c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) UpdatePreferences(c echo.Context) error {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return APIError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid user ID",
		}
	}

	var req account.UpdatePreferencesRequest
	if err := c.Bind(&req); err != nil {
		return APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Err:     err,
		}
	}

	updatedAccount, err := h.accountService.UpdatePreferences(c.Request().Context(), userID, req)
	if err != nil {
		return APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update preferences",
			Err:     err,
		}
	}

	return c.JSON(http.StatusOK, updatedAccount)
}

func (h *AccountHandler) Register(c echo.Context) error {
	var req account.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Err:     err,
		}
	}

	userAccount, err := h.accountService.Register(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrPasswordTooShort):
			return APIError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
				Err:     err,
			}
		case errors.Is(err, account.ErrUserAlreadyExists):
			return APIError{
				Code:    http.StatusConflict,
				Message: err.Error(),
				Err:     err,
			}
		default:
			return APIError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to register account",
				Err:     err,
			}
		}
	}

	response := map[string]any{
		"account": userAccount,
	}

	return c.JSON(http.StatusCreated, response)
}
