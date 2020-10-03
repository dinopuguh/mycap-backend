package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SigningKey is a secret key to store in jwt claims
var SigningKey = []byte(os.Getenv("MYCAP_JWT_TOKEN"))

// GenerateJWT creates JWT token from payload
func GenerateJWT(name, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	claims["email"] = email
	claims["issued"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(SigningKey)
	if err != nil {
		return "", fmt.Errorf("Failed to generate JWT")
	}

	return t, nil
}
