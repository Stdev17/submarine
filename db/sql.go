package db

import (
    "github.com/labstack/echo"

    //"encoding/json"
    "time"
    "log"
    "net/http"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

type Review struct {
    ReviewerID string `json:"reviewer"`
    Contents string `json:"contents"`
}

func CreateTable (c echo.Context) error {
    fmt.Println("Fuck")
    db, err := sql.Open("mysql", "root:$123@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }
    fmt.Println("Shit")
    err = db.Ping()
    if err != nil {
        return c.String(http.StatusInternalServerError, "db went wrong")
    }
    fmt.Println("You")
    fmt.Println("Fuck")
    _, err = db.Exec("use testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "choosing went wrong")
    }
    fmt.Println("You")
    reviews, err := db.Prepare("create table reviews(id int NOT NULL AUTO_INCREMENT, reviewer varchar(30), time datetime, contents varchar(1000), PRIMARY KEY (id));")
    if err != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    fmt.Println("Fuck")
    _, err = reviews.Exec()
    if err != nil {
        return c.String(http.StatusInternalServerError, "exec went wrong")
    }
    fmt.Println("You")

    defer db.Close()
    return c.String(http.StatusOK, "Table created")


}

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
    in, errinsert := db.Prepare("INSERT INTO reviews VALUES(?, ?, ?, ?)")
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
        fmt.Println(errid)
        return c.String(http.StatusInternalServerError, "auto scanning went wrong")
    }

    t := time.Now()
    t.Format("2006-01-02 15:04:05")

    if err != nil {
        return c.String(http.StatusInternalServerError, "parse went wrong")
    }
    _, errinto := in.Exec(id, rev.ReviewerID, t, rev.Contents)
    if errinto != nil {
        return c.String(http.StatusInternalServerError, "insert went wrong")
    }

    return c.String(http.StatusOK, "Successfully added the review")
}