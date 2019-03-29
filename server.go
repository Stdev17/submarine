package main

import (
    "github.com/labstack/echo"

    "github.com/submarine/handler"
    "github.com/submarine/db"

    "fmt"
)

func main () {
    fmt.Println("Hello, Server")

    e := echo.New()

    //groups

    //middlewares and groups

    //routing
    e.GET("/", handler.MainPage())
    e.GET("/initiate", db.CreateTable)
    e.POST("/create", db.CreateSQL)

    e.Start(":8000")
}