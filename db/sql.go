package db

import (
    "github.com/labstack/echo"

    //"encoding/json"
    "time"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

type Review struct {
    ReviewID int `json:"id"`
    ReviewerID string `json:"reviewer"`
    Time time.Time `json:"time"`
    Contents string `json:"contents"`
}

func CreateTable () error {
    db, err := sql.Open("mysql", "user:password@/maindb")
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }
    defer db.Close()

    //use DB
    _, err = db.Exec("create database reviewdb")
    if err != nil {
        return err
    }

    _, err = db.Exec("choose reviewdb")
    if err != nil {
        return err
    }

    reviews, err := db.Prepare("create table reviews(id int NOT NULL AUTO_INCREMENT, reviewer varchar(30), time unix_timestamp, contents varchar(1000), PRIMARY KEY (id));")
    if err != nil {
        return err
    }

    _, err = reviews.Exec()
    if err != nil {
        return err
    }

    return nil
}

func CreateSQL (c echo.Context) error {
    db, err := sql.Open("mysql", "user:password@/maindb")
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }
    defer db.Close()

    rev := Review{}

	defer c.Request().Body.Close()
    
    revjson := c.Bind(&rev)

    if revjson != nil {
        log.Printf("Failed adding a review: %s", err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    //use DB
    in, errinsert := db.Prepare("INSERT INTO reviews VALUES(?, ?, ?, ?)")
    if errinsert != nil {
        return err
    }
    defer in.Close()

    auto, errauto := db.Query("SELECT AUTO_INCREMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'reviewdb' and TABLE_NAME = 'reviews'")
    if errauto != nil {
        return err
    }
    defer auto.Close()

    var id int

    errid := auto.Scan(&id)
    if errid != nil {
        return errid
    }

    _, errinto := in.Exec(id, rev.ReviewerID, time.Now().Unix(), rev.Contents)
    if errinto != nil {
        return errinto
    }

    return c.String(http.StatusOK, "Successfully added the review")
}