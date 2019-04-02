package mymiddleware

import (
    "github.com/labstack/echo"
	//"log"
	//"net/http"
)

func JWTHeader (next echo.HandlerFunc) echo.HandlerFunc {
    return func (c echo.Context) error {
		cookie, _ := c.Cookie("login")

    	c.Request().Header.Set(echo.HeaderAuthorization, cookie.Value)

        return next(c)
    }
}