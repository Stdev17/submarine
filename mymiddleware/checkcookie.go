package mymiddleware

import (
    "github.com/labstack/echo"
    "net/http"
    "strings"
    "log"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func checkCookie (next echo.HandlerFunc) echo.HandlerFunc {
    return func (c echo.Context) error {

        cookie, err := c.Cookie("login")
        
        if err != nil {
            if strings.Contains(err.Error(), "named cookie not present") {
                return c.String(http.StatusUnauthorized, "you don't have any cookie")
            }
            log.Println(err)
            return err
        }

        check, err := ContainsValue(c, cookie.Value)
        if err != nil {
            log.Println(err)
            return err
        }
        if check {
            return next(c)
        }

        return c.String(http.StatusUnauthorized, "you don't have the right cookie")
    }
}

func ContainsValue (c echo.Context, value string) (bool, error) {
    db, err := sql.Open("mysql", "root:$123@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return false, err
    }

    err = db.Ping()
    if err != nil {
        return false, err
    }
    defer db.Close()

	defer c.Request().Body.Close()

    //use DB
    out, errOut := db.Prepare("select count(*) from users where jwt = ?;")
    if errOut != nil {
        return false, err
    }
    defer out.Close()

    auto, errRes := out.Query(value)
    if errRes != nil {
        return false, errRes
    }
    defer auto.Close()

    var chk int

    for auto.Next() {
        err := auto.Scan(&chk)
        if err != nil {
            return false, err
        }
    }

    errAuto := auto.Err()
    if errAuto != nil {
        return false, errAuto
    }

    if chk != 1 {
        return false, c.String(http.StatusInternalServerError, "Not unique ids")
    }

    return true, nil
}