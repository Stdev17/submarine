package main

import (
    "github.com/labstack/echo"
    "github.com/labstack/echo/engine/standard"
    "github.com/labstack/echo/middleware"

    "./handler"
    "./db"
)

func main () {

    e := echo.New()

    //groups

    //middlewares and groups
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    //routing
    e.Get("/main", hello.MainPage())
    e.Get("/main", sql.CreateTable())
    e.Post("/create", sql.CreateSQL())

    e.Run(standard.New(":8000"))
}