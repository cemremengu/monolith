package api

import (
	"errors"
	"net/http"

	"monolith/internal/service/account"
	"monolith/internal/service/auth"

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
	user, ok := c.Get("user").(*auth.AuthUser)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	account, err := h.accountService.GetAccountByID(c.Request().Context(), user.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Account not found").SetInternal(err)
	}

	return c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) UpdatePreferences(c echo.Context) error {
	user, ok := c.Get("user").(*auth.AuthUser)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
	}

	var req account.UpdatePreferencesRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").SetInternal(err)
	}

	updatedAccount, err := h.accountService.UpdatePreferences(c.Request().Context(), user.UserID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update preferences").SetInternal(err)
	}

	return c.JSON(http.StatusOK, updatedAccount)
}

func (h *AccountHandler) Register(c echo.Context) error {
	var req account.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").SetInternal(err)
	}

	userAccount, err := h.accountService.Register(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrPasswordTooShort):
			return echo.NewHTTPError(http.StatusBadRequest, "Password too short").SetInternal(err)
		case errors.Is(err, account.ErrUserAlreadyExists):
			return echo.NewHTTPError(http.StatusConflict, "User already exists").SetInternal(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to register user").SetInternal(err)
		}
	}

	response := map[string]any{
		"account": userAccount,
	}

	return c.JSON(http.StatusCreated, response)
}
