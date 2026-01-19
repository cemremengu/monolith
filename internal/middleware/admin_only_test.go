package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"monolith/internal/service/auth"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminOnly(t *testing.T) {
	tests := []struct {
		name         string
		user         *auth.AuthUser
		wantStatus   int
		handlerCalls int
	}{
		{
			name: "admin user",
			user: &auth.AuthUser{
				AccountID: uuid.New(),
				Email:     "admin@example.com",
				IsAdmin:   true,
				SessionID: uuid.New(),
			},
			wantStatus:   http.StatusOK,
			handlerCalls: 1,
		},
		{
			name: "non-admin user",
			user: &auth.AuthUser{
				AccountID: uuid.New(),
				Email:     "user@example.com",
				IsAdmin:   false,
				SessionID: uuid.New(),
			},
			wantStatus:   http.StatusNotFound,
			handlerCalls: 0,
		},
		{
			name:         "no user in context",
			user:         nil,
			wantStatus:   http.StatusNotFound,
			handlerCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.user != nil {
				c.Set("user", tt.user)
			}

			handlerCalled := 0
			handler := func(c echo.Context) error {
				handlerCalled++
				return c.String(http.StatusOK, "OK")
			}

			middleware := AdminOnly()
			err := middleware(handler)(c)

			if tt.wantStatus == http.StatusOK {
				require.NoError(t, err)
				assert.Equal(t, tt.handlerCalls, handlerCalled)
			} else {
				require.Error(t, err)
				httpErr := &echo.HTTPError{}
				ok := errors.As(err, &httpErr)
				require.True(t, ok)
				assert.Equal(t, tt.wantStatus, httpErr.Code)
			}
		})
	}
}

func TestAdminOnly_WrongUserType(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user", "not an AuthUser")

	handlerCalled := false
	handler := func(c echo.Context) error {
		handlerCalled = true
		return c.String(http.StatusOK, "OK")
	}

	middleware := AdminOnly()
	err := middleware(handler)(c)

	require.Error(t, err)
	httpErr := &echo.HTTPError{}
	ok := errors.As(err, &httpErr)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
	assert.False(t, handlerCalled)
}
