package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// SigningKey is a secret key to store in jwt claims
var SigningKey = []byte(os.Getenv("MYCAP_JWT_TOKEN"))

// JwtCustomClaims stores custom MyCap payload (name, email)
type JwtCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT creates JWT token from payload
func GenerateJWT(name, email string) (string, error) {
	claims := &JwtCustomClaims{
		name,
		email,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(SigningKey)
	if err != nil {
		return "", fmt.Errorf("Failed to generate JWT")
	}

	return t, nil
}

// GetJWTClaims gets user's payload
func GetJWTClaims(c echo.Context) *JwtCustomClaims {
	u := c.Get("user").(*jwt.Token)
	return u.Claims.(*JwtCustomClaims)
}
