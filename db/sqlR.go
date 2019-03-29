package db

import (
    "github.com/labstack/echo"

    "time"
    "net/http"
    "strconv"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

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
    out, errOut := db.Prepare("select * from reviews order by id desc limit 5 offset ?;")
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