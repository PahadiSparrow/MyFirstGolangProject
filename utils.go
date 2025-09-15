package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("PahadiBro") // Change this

// Generate JWT Token
func GenerateJWT(userID int, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = userID // Add user ID here
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix() // Token valid for 72 hours

	return token.SignedString(jwtSecret)
}

// Hash password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// Check password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
