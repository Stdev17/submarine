package handler

import (
	"github.com/labstack/echo"
	"net/http"
)
func MainJWT (c echo.Context) error {
	return c.String(http.StatusOK, "Top Secret")
}