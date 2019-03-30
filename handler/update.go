package handler

import (
    "github.com/labstack/echo"
    "net/http"
    "time"
    "log"
    "github.com/submarine/db"
	"github.com/submarine/config"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func Update (c echo.Context) error {
    data, err := sql.Open("mysql", "root:"+config.key.DB+"@tcp(127.0.0.1:3306)/testdb")
    if err != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }

    err = data.Ping()
    if err != nil {
        return c.String(http.StatusInternalServerError, "ping went wrong")
    }
    defer data.Close()

    user := db.user{}
    defer c.Request().Body.Close()    
    userjson := c.Bind(&user)

    if userjson != nil {
        log.Printf("Failed loading a user info: %s", err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    //use data
    cookie, err := c.Cookie("login")
    if err != nil {
        return c.String(http.StatusInternalServerError, "cookie went wrong")
    }

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
    	if _, ok := token.Method.(*jwt.SigningMethodHS512); !ok {
    	    return nil, c.String(http.StatusInternalServerError, "cookie went wrong")
    	}

    	return config.key.JWT, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	    //ok
	} else {
    	return c.String(http.StatusInternalServerError, "cookie went wrong")
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