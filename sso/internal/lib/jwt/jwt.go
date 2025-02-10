package jwt

import (
	"sso/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.User, app models.App, ttl time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"phone": user.Phone,
		"exp":   time.Now().Add(ttl).Unix(),
		"app":   app.Id,
	})

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
