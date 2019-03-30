package handler

import (
	"github.com/labstack/echo"
	"net/http"
)
func mainJWT (c echo.Context) error {
	return c.String(http.StatusOK, "Top Secret")
}