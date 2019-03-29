package db

import "time"

type Review struct {
    ReviewID int `json:"id"`
    ReviewerID string `json:"reviewer"`
    Time time.Time `json:"time"`
    Contents string `json:"contents"`
}