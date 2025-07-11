package api

import (
	"context"
	"log/slog"
	"net/http"

	"monolith/internal/config"
	"monolith/internal/database"
	customMiddleware "monolith/internal/middleware"
	"monolith/internal/service/account"
	"monolith/internal/service/auth"
	"monolith/internal/service/session"
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
	authService    *auth.Service
	sessionService *session.Service
}

// NewHTTPServer creates a new server instance with the given database, logger, and services.
func NewHTTPServer(
	db *database.DB,
	log *slog.Logger,
	cfg *config.Config,
	userService *user.Service,
	accountService *account.Service,
	authService *auth.Service,
	sessionService *session.Service,
) *HTTPServer {
	return &HTTPServer{
		echo:           echo.New(),
		db:             db,
		log:            log,
		config:         cfg,
		userService:    userService,
		accountService: accountService,
		authService:    authService,
		sessionService: sessionService,
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

	e.Use(customMiddleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())
	e.Use(middleware.RequestID())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(web.Assets()),
		HTML5:      true,
	}))

	hs.setupRoutes()
}

// setupRoutes configures all the application routes.
func (hs *HTTPServer) setupRoutes() {
	userHandler := NewUserHandler(hs.userService)
	authHandler := NewAuthHandler(hs.authService, hs.sessionService)
	accountHandler := NewAccountHandler(hs.authService, hs.accountService, hs.sessionService)

	api := hs.echo.Group("/api")

	// Public auth routes
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/logout", authHandler.Logout)

	// Protected routes
	securityConfig := config.NewSecurityConfig()
	protected := api.Group("", customMiddleware.SessionAuth(hs.sessionService, hs.accountService, securityConfig))

	protected.POST("/auth/rotate-token", authHandler.RotateToken)

	// Account management (profile & preferences)
	protected.GET("/account/profile", accountHandler.Profile)
	protected.PATCH("/account/preferences", accountHandler.UpdatePreferences)
	protected.GET("/account/sessions", accountHandler.GetSessions)
	protected.DELETE("/account/sessions/:sessionId", accountHandler.RevokeSession)
	// User administration (admin only)
	protected.GET("/users", userHandler.GetUsers)
	protected.GET("/users/:id", userHandler.GetUser)
	protected.POST("/users", userHandler.CreateUser)
	protected.PUT("/users/:id", userHandler.UpdateUser)
	protected.DELETE("/users/:id", userHandler.DeleteUser)
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
