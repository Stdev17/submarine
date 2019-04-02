package main

import (
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"

	//jwt "github.com/dgrijalva/jwt-go"

    "github.com/submarine/handler"
    "github.com/submarine/mymiddleware"
    "github.com/submarine/db"
    "github.com/submarine/config"
)

func main () {
    e := echo.New()

    //groups
    UpdateGroup := e.Group("/update")
    JWTGroup := e.Group("/jwt")

    //middlewares and groups
    e.Use(mymiddleware.ServerHeader)
    UpdateGroup.Use(mymiddleware.CheckCookie)
	//UpdateGroup.Use(mymiddleware.JWTHeader)
	config := middleware.JWTConfig{
            Claims:     &handler.JWTClaims{},
			SigningMethod: "HS512",
            SigningKey: []byte(config.Key.JWT),
    }
	UpdateGroup.Use(middleware.JWTWithConfig(config))
	UpdateGroup.POST("", handler.Update)
    JWTGroup.GET("", handler.MainJWT)

    //routing
    e.GET("/", handler.MainPage())
    e.GET("/login", handler.Login)

    e.GET("/initiate", db.CreateTable)
    e.GET("/user", db.CreateUserTable)
    e.POST("/create", db.CreateSQL)
    e.GET("/read", db.ReadSQL)

    e.POST("/register", db.CreateUser)

    e.Start(":8000")
}