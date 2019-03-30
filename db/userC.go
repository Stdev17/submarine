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
    db, err := sql.Open("mysql", "root:"+config.Key.DB+"@tcp(127.0.0.1:3306)/testdb")
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
    in, errinsert := db.Prepare("INSERT INTO users (userid, registered_time, hash) VALUES(?, ?, ?);")
    if errinsert != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    defer in.Close()

	chk, errchk := db.Prepare("select count(*) from users where userid = ?;")
    if errchk != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }
	defer chk.Close()
	var taken int
	err = chk.QueryRow(user.UserID).Scan(&taken)
	if err != nil {
		return c.String(http.StatusInternalServerError, "scanning went wrong")
	}
    if taken != 0 {
		return c.String(http.StatusBadRequest, "id already taken")
	}

    t := time.Now()

    if err != nil {
        return c.String(http.StatusInternalServerError, "parse went wrong")
    }

	hash := []byte(user.Hash)
	hash, err = bcrypt.GenerateFromPassword(hash, 8)
	if err != nil {
		log.Printf("%s", err)
        return c.String(http.StatusInternalServerError, "Hash went wrong")
    }
    _, errinto := in.Exec(user.UserID, t, hash)
    if errinto != nil {
        return c.String(http.StatusInternalServerError, "insert went wrong")
    }

    return c.String(http.StatusOK, "Successfully added the user")
}