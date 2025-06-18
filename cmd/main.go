package main

import (
	"log"
	"os"

	"monolith/internal/database"
	"monolith/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	userHandler := handlers.NewUserHandler(db)

	api := e.Group("/api")
	api.GET("/users", userHandler.GetUsers)
	api.GET("/users/:id", userHandler.GetUser)
	api.POST("/users", userHandler.CreateUser)
	api.PUT("/users/:id", userHandler.UpdateUser)
	api.DELETE("/users/:id", userHandler.DeleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(e.Start(":" + port))
}
