package handler_test

import (
	"testing"
    "github.com/stretchr/testify/assert"
	
	"github.com/labstack/echo"
	"github.com/submarine/handler"

	"net/http/httptest"
	"net/http"
	"net/url"
	//"log"
	
)

func TestLogin (t *testing.T) {
	assert := assert.New(t)

	e := echo.New()
	f := echo.New()
	q := make(url.Values)
	q.Set("username", "Jackd")
	q.Set("password", "1234")
	r := make(url.Values)
	r.Set("username", "Jack")
	r.Set("password", "1234")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	reqd := httptest.NewRequest(http.MethodGet, "/?"+r.Encode(), nil)
	recd := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	d := f.NewContext(reqd, recd)
	c.SetPath("/login")
	d.SetPath("/login")

	handler.Login(c)
	handler.Login(d)

	assert.NotEqual(rec.Code, recd.Code, "Login Test")

}