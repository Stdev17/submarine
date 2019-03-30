package handler

import (
    "github.com/labstack/echo"
    "net/http"
    "time"
    "log"
    "github.com/submarine/db"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func Update (c echo.Context) error {
    data, err := sql.Open("mysql", "root:$123@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    err = data.Ping()
    if err != nil {
        return c.String(http.StatusInternalServerError, "ping went wrong")
    }
    defer data.Close()

    rev := db.Review{}
    defer c.Request().Body.Close()    
    revjson := c.Bind(&rev)

    if revjson != nil {
        log.Printf("Failed updating a review: %s", err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    //use data
    out, errOut := data.Prepare("select jwt from users where userid = ?;")
    if errOut != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    defer out.Close()

    auto, errRes := out.Query(rev.ReviewerID)
    if errRes != nil {
        return c.String(http.StatusInternalServerError, "getting result went wrong")
    }
    defer auto.Close()

    var token string

    for auto.Next() {
        err := auto.Scan(&token)
        if err != nil {
            return c.String(http.StatusInternalServerError, "scanning went wrong")
        }
    }

    errAuto := auto.Err()
    if errAuto != nil {
        return c.String(http.StatusInternalServerError, "auto went wrong")
    }

    cookie, err := c.Cookie("login")
    if err != nil {
        return c.String(http.StatusInternalServerError, "cookie went wrong")
    }
    if token != cookie.Value {
        return c.String(http.StatusInternalServerError, "wrong cookie")
    }

    in, errup := data.Prepare("update reviews set contents = ? latest_time = ? where id = ? and reviewer = ?")
    if errup != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    defer in.Close()

    t := time.Now()

    if err != nil {
        return c.String(http.StatusInternalServerError, "parse went wrong")
    }
    _, errinto := in.Exec(rev.Contents, t, rev.ReviewID, rev.ReviewerID)
    if errinto != nil {
        return c.String(http.StatusInternalServerError, "insert went wrong")
    }

    return c.String(http.StatusOK, "Successfully updated the review")
}