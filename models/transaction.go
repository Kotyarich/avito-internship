package models

import "time"

type Transaction struct {
	UserId   int64     `json:"user_id"`
	Amount   float32   `json:"amount"`
	TargetId int64     `json:"target_id"`
	Type     string    `json:"type"`
	Time     time.Time `json:"time"`
}
