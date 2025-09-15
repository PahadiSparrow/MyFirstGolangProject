package main

import (
	"net/http"
	"strconv"

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
	var userID int
	err := DB.QueryRow("SELECT id, password FROM users WHERE email=?", loginReq.Email).Scan(&userID, &storedPassword)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	if !CheckPasswordHash(loginReq.Password, storedPassword) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	token, err := GenerateJWT(userID, loginReq.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Token generation failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func SubscribeCourse(c echo.Context) error {
	userID := c.Get("user_id").(int) // Extract user_id from context

	req := struct {
		CourseID int `json:"course_id"`
	}{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	stmt, _ := DB.Prepare("INSERT INTO subscriptions (user_id, course_id, subscribed_at) VALUES (?, ?, NOW())")
	_, err := stmt.Exec(userID, req.CourseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to subscribe"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Subscribed successfully"})
}

// Handler to get all users by course_id
func GetUsersByCourse(c echo.Context) error {
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid course_id"})
	}

	query := `
        SELECT u.name, u.email, u.phone
        FROM users u
        JOIN subscriptions cs ON u.id = cs.user_id
        WHERE cs.course_id = ?`

	rows, err := DB.Query(query, courseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	users := []map[string]string{}
	for rows.Next() {
		var name, email, phone string
		if err := rows.Scan(&name, &email, &phone); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		users = append(users, map[string]string{
			"name":  name,
			"email": email,
			"phone": phone,
		})
	}

	return c.JSON(http.StatusOK, users)
}
