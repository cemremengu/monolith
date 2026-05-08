package api

import (
	"errors"
	"net/http"

	"monolith/internal/service/account"
	"monolith/internal/service/auth"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type AccountHandler struct {
	accountService *account.Service
}

func NewAccountHandler(accountService *account.Service) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

func (h *AccountHandler) Profile(c *echo.Context) error {
	user, ok := c.Get("user").(*auth.AuthUser)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user")
	}

	userID := user.AccountID

	account, err := h.accountService.GetAccountByID(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Account not found").Wrap(err)
	}

	return c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) UpdatePreferences(c *echo.Context) error {
	user, ok := c.Get("user").(*auth.AuthUser)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user")
	}

	userID := user.AccountID

	var req account.UpdatePreferencesRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").Wrap(err)
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	updatedAccount, err := h.accountService.UpdatePreferences(c.Request().Context(), userID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update preferences").Wrap(err)
	}

	return c.JSON(http.StatusOK, updatedAccount)
}

func (h *AccountHandler) Register(c *echo.Context) error {
	var req account.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").Wrap(err)
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	userAccount, err := h.accountService.Register(c.Request().Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrPasswordTooShort):
			return echo.NewHTTPError(http.StatusBadRequest, "Password too short").Wrap(err)
		case errors.Is(err, account.ErrUserAlreadyExists):
			return echo.NewHTTPError(http.StatusConflict, "User already exists").Wrap(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to register user").Wrap(err)
		}
	}

	response := map[string]any{
		"account": userAccount,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *AccountHandler) GetAccounts(c *echo.Context) error {
	accounts, err := h.accountService.GetAccounts(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve accounts").Wrap(err)
	}

	return c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) GetAccount(c *echo.Context) error {
	id := c.Param("id")
	accountID, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid account ID format").Wrap(err)
	}

	account, err := h.accountService.GetAccount(c.Request().Context(), accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Account not found").Wrap(err)
	}

	return c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) CreateAccount(c *echo.Context) error {
	var req account.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").Wrap(err)
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	createdAccount, err := h.accountService.CreateAccount(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create account").Wrap(err)
	}

	return c.JSON(http.StatusCreated, createdAccount)
}

func (h *AccountHandler) UpdateAccount(c *echo.Context) error {
	id := c.Param("id")
	accountID, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid account ID format").Wrap(err)
	}

	var req account.UpdateAccountRequest
	if err = c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").Wrap(err)
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	updatedAccount, err := h.accountService.UpdateAccount(c.Request().Context(), accountID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update account").Wrap(err)
	}

	return c.JSON(http.StatusOK, updatedAccount)
}

func (h *AccountHandler) DisableAccount(c *echo.Context) error {
	id := c.Param("id")
	accountID, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid account ID format").Wrap(err)
	}

	err = h.accountService.DisableAccount(c.Request().Context(), accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to disable account").Wrap(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AccountHandler) EnableAccount(c *echo.Context) error {
	id := c.Param("id")
	accountID, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid account ID format").Wrap(err)
	}

	err = h.accountService.EnableAccount(c.Request().Context(), accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to enable account").Wrap(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AccountHandler) DeleteAccount(c *echo.Context) error {
	id := c.Param("id")
	accountID, err := uuid.Parse(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid account ID format").Wrap(err)
	}

	err = h.accountService.DeleteAccount(c.Request().Context(), accountID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete account").Wrap(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AccountHandler) InviteUsers(c *echo.Context) error {
	var req account.InviteUsersRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").Wrap(err)
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	response, err := h.accountService.InviteUsers(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to invite users").Wrap(err)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AccountHandler) ChangePassword(c *echo.Context) error {
	user, ok := c.Get("user").(*auth.AuthUser)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user")
	}

	var req account.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body").Wrap(err)
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	err := h.accountService.ChangePassword(c.Request().Context(), user.AccountID, req)
	if err != nil {
		switch {
		case errors.Is(err, account.ErrPasswordTooShort):
			return echo.NewHTTPError(http.StatusBadRequest, "Password too short").Wrap(err)
		case errors.Is(err, account.ErrInvalidPassword):
			return echo.NewHTTPError(http.StatusBadRequest, "Current password is incorrect").Wrap(err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to change password").Wrap(err)
		}
	}

	return c.NoContent(http.StatusNoContent)
}
