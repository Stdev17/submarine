package mymiddleware

import (
    "github.com/labstack/echo"
	"github.com/submarine/config"
    "net/http"
    "strings"
    "log"

	jwt "github.com/dgrijalva/jwt-go"
)

func CheckCookie (next echo.HandlerFunc) echo.HandlerFunc {
    return func (c echo.Context) error {

        cookie, err := c.Cookie("login")
        
        if err != nil {
            if strings.Contains(err.Error(), "named cookie not present") {
                return c.String(http.StatusUnauthorized, "you don't have any cookie")
            }
            log.Println(err)
            return err
        }

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				
				return nil, c.String(http.StatusInternalServerError, "cookie went wrong")
			}

			return config.Key.JWT, nil
		})

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return next(c)
		} else {
    		return c.String(http.StatusUnauthorized, "you don't have the right cookie")
		}
    }
}