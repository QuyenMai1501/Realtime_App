package auth

import (
	"WS_GIN_GOZIL/src/common"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(common.GetEnv("JWT_SECRET"))

func GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}