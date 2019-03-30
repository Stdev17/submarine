package db

import (
    "github.com/labstack/echo"

    "time"
    "log"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func CreateSQL (c echo.Context) error {
    db, err := sql.Open("mysql", "root:$123@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    err = db.Ping()
    if err != nil {
        return c.String(http.StatusInternalServerError, "ping went wrong")
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
    in, errinsert := db.Prepare("INSERT INTO reviews (reviewer, time, latest_time, contents) VALUES(?, ?, ?, ?)")
    if errinsert != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    defer in.Close()

    auto, errauto := db.Query("SELECT AUTO_INCREMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'testdb' and TABLE_NAME = 'reviews';")
    if errauto != nil {
        return c.String(http.StatusInternalServerError, "taking auto went wrong")
    }
    defer auto.Close()

    var id int

    auto.Next()
    errid := auto.Scan(&id)
    if errid != nil {
        return c.String(http.StatusInternalServerError, "auto scanning went wrong")
    }

    t := time.Now()

    if err != nil {
        return c.String(http.StatusInternalServerError, "parse went wrong")
    }
    _, errinto := in.Exec(rev.ReviewerID, t, t, rev.Contents)
    if errinto != nil {
        return c.String(http.StatusInternalServerError, "insert went wrong")
    }

    return c.String(http.StatusOK, "Successfully added the review")
}