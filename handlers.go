package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Register User Endpoint
func RegisterUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	stmt, err := DB.Prepare("INSERT INTO users (name, email, phone, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "DB error"})
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Name, u.Email, u.Phone, hashedPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User already exists or DB error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

// Login User Endpoint
func LoginUser(c echo.Context) error {
	loginReq := new(User)
	if err := c.Bind(loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	var storedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE email=?", loginReq.Email).Scan(&storedPassword)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	if !CheckPasswordHash(loginReq.Password, storedPassword) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	token, err := GenerateJWT(loginReq.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Token generation failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
