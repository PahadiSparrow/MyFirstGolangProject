package main

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// JWT Middleware to protect routes
func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
		}

		tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", int(claims["user_id"].(float64))) // Save user_id in context

		return next(c)
	}
}
