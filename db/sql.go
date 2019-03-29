package db

import (
    "github.com/labstack/echo"

    //"encoding/json"
    "time"
    "log"
    "net/http"
    "fmt"
    "strconv"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

type Review struct {
    ReviewID int `json:"id"`
    ReviewerID string `json:"reviewer"`
    Time time.Time `json:"time"`
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
    reviews, err := db.Prepare("create table reviews(id int NOT NULL AUTO_INCREMENT, reviewer varchar(30), time timestamp, contents varchar(1000), PRIMARY KEY (id)) engine=innodb;")
    if err != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    fmt.Println("Fuck")
    _, err = reviews.Exec()
    if err != nil {
        return c.String(http.StatusInternalServerError, "exec went wrong")
    }

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
    in, errinsert := db.Prepare("INSERT INTO reviews (reviewer, time, contents) VALUES(?, ?, ?)")
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

    if err != nil {
        return c.String(http.StatusInternalServerError, "parse went wrong")
    }
    _, errinto := in.Exec(rev.ReviewerID, t, rev.Contents)
    if errinto != nil {
        fmt.Printf("%d %s %s %s", id, rev.ReviewerID, t, rev.Contents)
        return c.String(http.StatusInternalServerError, "insert went wrong")
    }

    return c.String(http.StatusOK, "Successfully added the review")
}

func ReadSQL (c echo.Context) error {
    db, err := sql.Open("mysql", "root:$123@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    err = db.Ping()
    if err != nil {
        return c.String(http.StatusInternalServerError, "ping went wrong")
    }
    defer db.Close()

    var offset int
    offset, _ = strconv.Atoi(c.QueryParam("offset"))

	defer c.Request().Body.Close()

    //use DB
    out, errOut := db.Prepare("select * from reviews orders limit 5 offset ?;")
    if errOut != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    defer out.Close()

    var (
        id int
        reviewer string
        myTime string
        contents string
    )

    auto, errRes := out.Query(offset)
    if errRes != nil {
        return c.String(http.StatusInternalServerError, "taking result went wrong")
    }
    defer auto.Close()

    var res []Review

    for auto.Next() {
        err := auto.Scan(&id, &reviewer, &myTime, &contents)
        if err != nil {
            return c.String(http.StatusInternalServerError, "taking rows went wrong")
        }
        timeStamp, err := time.Parse("2006-01-02 15:04:05", myTime)
        if err != nil {
            return c.String(http.StatusInternalServerError, "parsing times went wrong")
        }
        res = append(res, Review{
            ReviewID: id,
            ReviewerID: reviewer,
            Time: timeStamp,
            Contents: contents,
        })
    }
    errAuto := auto.Err()
    if errAuto != nil {
        return c.String(http.StatusInternalServerError, "taking rows went wrong")
    }

    if res != nil {
        return c.JSON(http.StatusOK, res)
    }

    return c.String(http.StatusInternalServerError, "no data")
}