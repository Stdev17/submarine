package handler

import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
    Name string `json:"name"`
    jwt.StandardClaims
}

func createJWTToken (id, password string) (string, error) {
    claims := JWTClaims{
        id,
        jwt.StandardClaims{
            Id: id,
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
    }

    rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

    token, err := rawToken.SignedString([]byte(password))
    if err != nil {
        return "", err
    }

    return token, err
}