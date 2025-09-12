package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public routes
	e.POST("/register", RegisterUser)
	e.POST("/login", LoginUser)

	// Start server
	log.Println("Starting server on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
