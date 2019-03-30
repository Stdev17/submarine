package handler

import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/submarine/config"
)

type JWTClaims struct {
    Name string `json:"name"`
    jwt.StandardClaims
}

func createJWTToken (id, password []byte) (string, error) {
    claims := JWTClaims{
        string(id),
        jwt.StandardClaims{
            Id: string(id),
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
    }

    rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

    token, err := rawToken.SignedString([]byte(config.Key.JWT))
    if err != nil {
        return "", err
    }

    return token, err
}