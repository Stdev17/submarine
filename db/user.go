package db

import "time"

type User struct {
    ID int `json:"id"`
    UserID string `json:"user"`
    Time time.Time `json:"registered"`
    Hash string `json:"hash"`
}