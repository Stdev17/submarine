package handler

import (
	"github.com/labstack/echo"
	"net/http"
	"time"
	"log"
)

func Login (c echo.Context) error {
    username := c.QueryParam("username")
    password := c.QueryParam("password")

    if username == "Jack" && password == "1234" {
        cookie := &http.Cookie{}

        cookie.Name = "sessionID"
        cookie.Value = "some_hash"
        cookie.Expires = time.Now().Add(48 * time.Hour)

        c.SetCookie(cookie)

        // create jwt token
        token, err := createJWTToken()
        if err != nil {
            log.Println("Error Creating JWT Tokens", err)
            return c.String(http.StatusInternalServerError, "something went wrong")
        }

        return c.JSON(http.StatusOK, map[string]string{
            "message": "You were logged in!",
            "token": token,
        })
    }

    return c.String(http.StatusUnauthorized, "Your username or password is invalid.")
}