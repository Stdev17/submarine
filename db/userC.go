package db

import (
    "github.com/labstack/echo"

    "time"
    "log"
    "net/http"

	"golang.org/x/crypto/bcrypt"
	"github.com/submarine/config"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func CreateUser (c echo.Context) error {
    db, err := sql.Open("mysql", "root:"+config.key.DB+"@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    err = db.Ping()
    if err != nil {
        return c.String(http.StatusInternalServerError, "ping went wrong")
    }
    defer db.Close()

    user := User{}
	defer c.Request().Body.Close()    
    userjson := c.Bind(&user)

    if userjson != nil {
        log.Printf("Failed adding a user: %s", err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    //use DB
    in, errinsert := db.Prepare("INSERT INTO users (userid, time, hash) VALUES(?, ?, ?)")
    if errinsert != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    defer in.Close()

    t := time.Now()

    if err != nil {
        return c.String(http.StatusInternalServerError, "parse went wrong")
    }

	hash, err := bcrypt.GenerateFromPassword(user.Hash, 72)
	if err != nil {
        return c.String(http.StatusInternalServerError, "Hash went wrong")
    }
    _, errinto := in.Exec(user.UserID, t, hash)
    if errinto != nil {
        return c.String(http.StatusInternalServerError, "insert went wrong")
    }

    return c.String(http.StatusOK, "Successfully added the review")
}