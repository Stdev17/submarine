package handler

import (
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
    Name string `json:"name"`
    jwt.StandardClaims
}

func createJWTToken() (string, error) {
    claims := JWTClaims{
        "Jack",
        jwt.StandardClaims{
            Id: "main_user_id",
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
    }

    rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

    token, err := rawToken.SignedString([]byte("mySecret"))
    if err != nil {
        return "", err
    }

    return token, err
}