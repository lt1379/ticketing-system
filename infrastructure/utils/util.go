package utils

import (
	"my-project/infrastructure/logger"
	"time"

	"github.com/golang-jwt/jwt"
)

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}

func GenerateToken(payload map[string]interface{}, secretKey string) (string, error) {
	var claims jwt.MapClaims = payload
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while generate token")
		return "", err
	}
	return tokenString, nil
}
