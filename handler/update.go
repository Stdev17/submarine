package handler

import (
	"github.com/labstack/echo"
	//"strconv"
	"net/http"
)

func Update (c echo.Context) error {
	//reviewID, _ := strconv.Atoi(c.QueryParam("id"))

	return c.String(http.StatusInternalServerError, "")
}