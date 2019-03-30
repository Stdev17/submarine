package mymiddleware

import (
    "github.com/labstack/echo"
)

func ServerHeader (next echo.HandlerFunc) echo.HandlerFunc {
    return func (c echo.Context) error {
        c.Response().Header().Set(echo.HeaderServer, "myServer")

        return next(c)
    }
}