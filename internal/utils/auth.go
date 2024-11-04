package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/OPTIC7409/tutor-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func ExtractUserIDFromToken(c *fiber.Ctx, db *gorm.DB) (int, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("missing Authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("invalid Authorization header format")
	}

	tokenString := parts[1]

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return 0, errors.New("invalid token signature")
		}
		return 0, errors.New("invalid token")
	}

	if !token.Valid {
		return 0, errors.New("invalid token")
	}

	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return 0, errors.New("token expired")
	}

	var user models.User
	result := db.Where("id = ? AND token = ?", claims.UserID, tokenString).First(&user)
	if result.Error != nil {
		return 0, errors.New("token not found or invalid")
	}

	return claims.UserID, nil
}
