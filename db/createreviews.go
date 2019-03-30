package db

import (
	"github.com/labstack/echo"

    "net/http"
	"github.com/submarine/config"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func CreateTable (c echo.Context) error {
    db, err := sql.Open("mysql", "root:"+config.key.DB+"@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    err = db.Ping()
    if err != nil {
        return c.String(http.StatusInternalServerError, "db went wrong")
    }

    _, err = db.Exec("use testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "choosing went wrong")
    }
    
    reviews, err := db.Prepare("create table reviews(id int NOT NULL AUTO_INCREMENT, reviewer varchar(30), time timestamp, latest_time timestamp, contents varchar(1000), PRIMARY KEY (id)) engine=innodb;")
    if err != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    
    _, err = reviews.Exec()
    if err != nil {
        return c.String(http.StatusInternalServerError, "exec went wrong")
    }

    defer db.Close()
    return c.String(http.StatusOK, "Table created")
}

func CreateUserTable (c echo.Context) error {
    db, err := sql.Open("mysql", "root:$123@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    err = db.Ping()
    if err != nil {
        return c.String(http.StatusInternalServerError, "db went wrong")
    }

    _, err = db.Exec("use testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "choosing went wrong")
    }
    
    reviews, err := db.Prepare("create table users(id int NOT NULL AUTO_INCREMENT, userid varchar(30), registrered_time timestamp, hash varchar(255), PRIMARY KEY (id)) engine=innodb;")
    if err != nil {
        return c.String(http.StatusInternalServerError, "query went wrong")
    }
    
    _, err = reviews.Exec()
    if err != nil {
        return c.String(http.StatusInternalServerError, "exec went wrong")
    }

    defer db.Close()
    return c.String(http.StatusOK, "Table created")
}