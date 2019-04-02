package handler_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"

    "github.com/submarine/handler"
    "github.com/submarine/config"
	"github.com/submarine/mymiddleware"

    "net/http/httptest"
    "log"
	"bytes"
)

func TstUpdate (t *testing.T) {
    assert := assert.New(t)
	myBody := `
		"id": 9,
		"reviewer": "Jack",
		"contents": "Fuck you!"
	`
	myBodyd := `
		"id": 9,
		"reviewer": "Jacku",
		"contents": "Fuck you!"
	`
	reqBody := bytes.NewBufferString(myBody)
	reqBodyd := bytes.NewBufferString(myBodyd)

    e := echo.New()
	f := echo.New()
    req := httptest.NewRequest(echo.POST, "/update?username=Jack&password=1234", reqBody)
    rec := httptest.NewRecorder()
    reqd := httptest.NewRequest(echo.POST, "/update?username=Jack&password=1234", reqBodyd)
    recd := httptest.NewRecorder()
    //c := e.NewContext(req, rec)
    //d := f.NewContext(reqd, recd)
    UpdateGroup := e.Group("/update")
	UpGroup := f.Group("/update")

    e.Use(mymiddleware.ServerHeader)
    UpdateGroup.Use(mymiddleware.CheckCookie)
	UpdateGroup.Use(mymiddleware.JWTHeader)
	f.Use(mymiddleware.ServerHeader)
    UpGroup.Use(mymiddleware.CheckCookie)
	UpGroup.Use(mymiddleware.JWTHeader)

	config := middleware.JWTConfig{
            Claims:     &handler.JWTClaims{},
            SigningKey: []byte(config.Key.JWT),
    }
	UpdateGroup.Use(middleware.JWTWithConfig(config))
	UpdateGroup.POST("", handler.Update)
	UpGroup.Use(middleware.JWTWithConfig(config))
	UpGroup.POST("", handler.Update)
	//c := e.NewContext(req, rec)
	//d := e.NewContext(reqd, recd)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	reqd.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	e.ServeHTTP(rec, req)
	f.ServeHTTP(recd, reqd)

    assert.NotEqual(rec.Code, recd.Code, "Update Test")

	log.Println(rec)
	log.Println(recd)
	log.Println(req.Header.Get(echo.HeaderAuthorization))
	
}