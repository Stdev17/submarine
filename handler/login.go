package handler

import (
    "github.com/labstack/echo"
    "net/http"
    "time"
    "log"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func Login (c echo.Context) error {
    username := c.QueryParam("username")
    password := c.QueryParam("password")

    // create jwt token
    token, err := createJWTToken(username, password)
    if err != nil {
        log.Println("Error Creating JWT Tokens", err)
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    check, err := checkUser(username, token, c)
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong in db")
    }
    if !check {
        return c.String(http.StatusUnauthorized, "Your username or password is invalid.")
    }
    
    cookie := &http.Cookie{}

    cookie.Name = "login"
    cookie.Value = token
    cookie.Expires = time.Now().Add(48 * time.Hour)

    c.SetCookie(cookie)   

    return c.String(http.StatusOK, "You were logged in!")
}

func checkUser (id, password string, c echo.Context) (bool, error) {
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
    out, errOut := db.Prepare("select count(*) from users where userid = ? and password = ?;")
    if errOut != nil {
        return false, errOut
    }
    defer out.Close()

    auto, errRes := out.Query(id, password)
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