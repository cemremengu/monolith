package api

import (
	"log/slog"
	"net/http"
	"os"

	"monolith/internal/database"
	customMiddleware "monolith/internal/middleware"
	"monolith/web"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server wraps the Echo server and provides methods for setup and startup.
type Server struct {
	echo *echo.Echo
	db   *database.DB
	log  *slog.Logger
}

// NewServer creates a new server instance with the given database and logger.
func NewServer(db *database.DB, log *slog.Logger) *Server {
	return &Server{
		echo: echo.New(),
		db:   db,
		log:  log,
	}
}

// Setup configures the server with middleware and routes.
func (s *Server) Setup() {
	e := s.echo

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
	e.Use(middleware.Secure())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: http.FS(web.Assets()),
		HTML5:      true,
	}))

	s.setupRoutes()
}

// setupRoutes configures all the application routes.
func (s *Server) setupRoutes() {
	userHandler := NewUserHandler(s.db)
	authHandler := NewAuthHandler(s.db)

	api := s.echo.Group("/api")

	// Public auth routes
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)
	api.POST("/auth/refresh", authHandler.RefreshToken)

	// Protected routes
	protected := api.Group("", customMiddleware.JWTAuth())
	protected.GET("/auth/me", authHandler.Me)
	protected.POST("/auth/logout", authHandler.Logout)
	protected.GET("/auth/sessions", authHandler.GetSessions)
	protected.DELETE("/auth/sessions/:sessionId", authHandler.RevokeSession)
	protected.POST("/auth/sessions/revoke-others", authHandler.RevokeAllOtherSessions)
	protected.GET("/users", userHandler.GetUsers)
	protected.GET("/users/:id", userHandler.GetUser)
	protected.POST("/users", userHandler.CreateUser)
	protected.PUT("/users/:id", userHandler.UpdateUser)
	protected.DELETE("/users/:id", userHandler.DeleteUser)
}

// Start starts the server on the specified port.
func (s *Server) Start() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	s.log.Info("Server starting", slog.String("port", port))
	return s.echo.Start(":" + port)
}

// Echo returns the underlying Echo instance for advanced configuration if needed.
func (s *Server) Echo() *echo.Echo {
	return s.echo
}
