package models

import "time"

type Session struct {
	Model
	Token   string     `json:"token"`
	UserID  uint64     `json:"userID"`
	Expires *time.Time `json:"expires"`
}
