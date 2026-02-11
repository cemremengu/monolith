package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"monolith/internal/config"
	"monolith/internal/database"
	mw "monolith/internal/middleware"
	"monolith/internal/service/account"
	"monolith/internal/service/auth"
	"monolith/internal/service/login"
	"monolith/web"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HTTPServer wraps the Echo server and provides methods for setup and startup.
type HTTPServer struct {
	echo           *echo.Echo
	db             *database.DB
	config         *config.Config
	accountService *account.Service
	loginService   *login.Service
	authService    *auth.Service
}

// NewHTTPServer creates a new server instance with the given database, logger, and services.
func NewHTTPServer(
	db *database.DB,
	cfg *config.Config,
	accountService *account.Service,
	loginService *login.Service,
	authService *auth.Service,
) *HTTPServer {
	return &HTTPServer{
		echo:           echo.New(),
		db:             db,
		config:         cfg,
		accountService: accountService,
		loginService:   loginService,
		authService:    authService,
	}
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// Setup configures the server with middleware and routes.
func (hs *HTTPServer) Setup() {
	e := hs.echo

	e.HideBanner = true
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Pre(middleware.RemoveTrailingSlash())

	// Middleware
	// the order of the middleware is important in most cases

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(_ echo.Context, err error, stack []byte) error {
			slog.Error("[PANIC RECOVER]", "error", err, "stack", string(stack))
			return err
		},
	}))

	e.Use(mw.Logger())
	e.Use(middleware.Gzip())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/api")
		},
		Filesystem: http.FS(web.Assets()),
		HTML5:      true,
	}))

	hs.RegisterRoutes()
}

// Start starts the server on the specified port.
func (hs *HTTPServer) Start() error {
	port := hs.config.Server.Port

	slog.Info("Server starting", "port", port)
	return hs.echo.Start(":" + port)
}

// Shutdown gracefully shuts down the server.
func (hs *HTTPServer) Shutdown(ctx context.Context) error {
	return hs.echo.Shutdown(ctx)
}
