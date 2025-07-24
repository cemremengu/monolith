package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"monolith"

	"monolith/internal/config"
	"monolith/internal/database"
	mw "monolith/internal/middleware"
	"monolith/internal/service/account"
	"monolith/internal/service/auth"
	"monolith/internal/service/login"
	"monolith/internal/service/user"
	"monolith/web"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// HTTPServer wraps the Echo server and provides methods for setup and startup.
type HTTPServer struct {
	echo           *echo.Echo
	db             *database.DB
	log            *slog.Logger
	config         *config.Config
	userService    *user.Service
	accountService *account.Service
	loginService   *login.Service
	authService    *auth.Service
}

// NewHTTPServer creates a new server instance with the given database, logger, and services.
func NewHTTPServer(
	db *database.DB,
	log *slog.Logger,
	cfg *config.Config,
	userService *user.Service,
	accountService *account.Service,
	loginService *login.Service,
	authService *auth.Service,
) *HTTPServer {
	return &HTTPServer{
		echo:           echo.New(),
		db:             db,
		log:            log,
		config:         cfg,
		userService:    userService,
		accountService: accountService,
		loginService:   loginService,
		authService:    authService,
	}
}

// Setup configures the server with middleware and routes.
func (hs *HTTPServer) Setup() {
	e := hs.echo

	e.HideBanner = true

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
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.RequestID())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/api")
		},
		Filesystem: http.FS(web.Assets()),
		HTML5:      true,
	}))

	hs.setupRoutes()
}

// setupRoutes configures all the application routes.
func (hs *HTTPServer) setupRoutes() {
	userHandler := NewUserHandler(hs.userService)
	authHandler := NewAuthHandler(hs.loginService, hs.authService)
	accountHandler := NewAccountHandler(hs.accountService)
	authSessionHandler := NewSessionHandler(hs.authService)

	api := hs.echo.Group("/api")

	// Public auth routes
	api.POST("/login", authHandler.Login)
	api.POST("/logout", authHandler.Logout)
	api.GET("/version", func(c echo.Context) error {
		return c.JSON(http.StatusOK, monolith.GetVersionInfo())
	})
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Protected routes
	protected := api.Group("", mw.SessionAuth(hs.authService, hs.accountService, hs.config.Security))

	protected.GET("/account/sessions", authSessionHandler.GetSessions)
	protected.DELETE("/account/sessions/:sessionId", authSessionHandler.RevokeSession)
	protected.POST("/account/sessions/rotate", authSessionHandler.RotateSession)

	protected.GET("/account/profile", accountHandler.Profile)
	protected.PATCH("/account/preferences", accountHandler.UpdatePreferences)
	protected.POST("/account/register", accountHandler.Register)

	// Admin-only routes
	admin := protected.Group("", mw.AdminOnly())
	admin.GET("/users", userHandler.GetUsers)
	admin.GET("/users/:id", userHandler.GetUser)
	admin.POST("/users", userHandler.CreateUser)
	admin.PUT("/users/:id", userHandler.UpdateUser)
	admin.DELETE("/users/:id", userHandler.DeleteUser)
}

// Start starts the server on the specified port.
func (hs *HTTPServer) Start() error {
	port := hs.config.Server.Port

	hs.log.Info("Server starting", slog.String("port", port))
	return hs.echo.Start(":" + port)
}

// Shutdown gracefully shuts down the server.
func (hs *HTTPServer) Shutdown(ctx context.Context) error {
	return hs.echo.Shutdown(ctx)
}

// Echo returns the underlying Echo instance for advanced configuration if needed.
func (hs *HTTPServer) Echo() *echo.Echo {
	return hs.echo
}
