package handler

import (
    "github.com/labstack/echo"
    "net/http"
    "time"
    "log"

    "golang.org/x/crypto/bcrypt"
    "github.com/submarine/config"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)



func Login (c echo.Context) error {
    var username, password string
    username = c.QueryParam("username")
    password = c.QueryParam("password")
    
    check, err := checkUser(username, password, c)
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong in db")
    }
    if !check {
        return c.String(http.StatusUnauthorized, "Your username or password is invalid.")
    }

    token, err := createJWTToken([]byte(username), []byte(password))
    if err != nil {
        log.Println("Error Creating JWT Tokens", err)
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    
    cookie := &http.Cookie{}

    cookie.Name = "login"
    cookie.Value = token
    cookie.Expires = time.Now().Add(48 * time.Hour)

    c.SetCookie(cookie)   

    return c.String(http.StatusOK, "You were logged in!")
}

func checkUser (id, password string, c echo.Context) (bool, error) {
    db, err := sql.Open("mysql", "root:"+config.Key.DB+"@tcp(127.0.0.1:3306)/testdb")
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
    out, errOut := db.Prepare("select hash from users where userid = ?;")
    if errOut != nil {
        return false, errOut
    }
    defer out.Close()

    auto, errRes := out.Query(id)
    if errRes != nil {
        return false, errRes
    }
    defer auto.Close()

    var hash string

    for auto.Next() {
        err := auto.Scan(&hash)
        if err != nil {
            return false, err
        }
    }
    errAuto := auto.Err()
    if errAuto != nil {
        return false, errAuto
    }
    if hash == "" {
        return false, c.String(http.StatusBadRequest, "ID unidentified")
    }

    errchk := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if errchk != nil {
        log.Println(errchk)
        return false, errchk
    }

    return true, nil
}