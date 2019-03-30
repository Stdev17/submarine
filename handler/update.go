package handler

import (
    "github.com/labstack/echo"
    "net/http"
    "time"
    "log"
    "github.com/submarine/db"
    "github.com/submarine/config"

    jwt "github.com/dgrijalva/jwt-go"

    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func Update (c echo.Context) error {
    data, err := sql.Open("mysql", "root:"+config.Key.DB+"@tcp(127.0.0.1:3306)/testdb")
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
        log.Printf("Failed loading a review info: %s", err)
        return echo.NewHTTPError(http.StatusInternalServerError)
    }

    var revid string
    cookie, err := c.Cookie("login")
    token, err := jwt.ParseWithClaims(cookie.Value, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, c.String(http.StatusInternalServerError, "cookie went wrong")
        }
        return config.Key.JWT, nil
    })
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        log.Printf("%v %v", claims.Id, claims.ExpiresAt)
        revid = claims.Id
    } else {
        log.Println(err)
    }

    if rev.ReviewerID != revid {
        return c.String(http.StatusUnauthorized, "you does not own the review")
    }

    //use data
    chk, errchk := data.Prepare("select count(*) from reviews where id = ? and reviewer = ?;")
    if errchk != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }
    defer chk.Close()
    out, errchk := chk.Query(rev.ReviewID, rev.ReviewerID)
    if errchk != nil {
        return c.String(http.StatusInternalServerError, "something went wrong")
    }
    defer out.Close()
    var taken int
    out.Next()
    err = out.Scan(&taken)
    if err != nil {
        return c.String(http.StatusInternalServerError, "scanning went wrong")
    }
    if out.Err() != nil {
        return c.String(http.StatusInternalServerError, "scanning went wrong")	
    }
    if taken != 1 {
        return c.String(http.StatusBadRequest, "review finding error")
    }

    in, errup := data.Prepare("update reviews set contents = ?, latest_time = ? where id = ? and reviewer = ?;")
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