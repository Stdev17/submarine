package main

import (
    "github.com/labstack/echo"

    "github.com/submarine/handler"
    "github.com/submarine/mymiddleware"
    "github.com/submarine/db"
)

func main () {
    e := echo.New()

    //groups
    UpdateGroup := e.Group("/update")
    JWTGroup := e.Group("/jwt")

    //middlewares and groups
    e.Use(mymiddleware.ServerHeader)
    UpdateGroup.Use(mymiddleware.checkCookie)
    UpdateGroup.POST("/update", handler.Update)
    JWTGroup.GET("/", handler.mainJWT)

    //routing
    e.GET("/", handler.MainPage())
    e.GET("/login", handler.Login)

    e.GET("/initiate", db.CreateUserTable)
    e.POST("/create", db.CreateSQL)
    e.GET("/read", db.ReadSQL)

    e.Start(":8000")
}