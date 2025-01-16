package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(name string) (string, error) {

	claims := jwt.MapClaims{
		"userID": name,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Токен действует 24 часа
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(os.Getenv("KEY"))
	return token.SignedString(key)
}
