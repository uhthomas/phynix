package models

import "time"

type Verification struct {
	Model
	Token    string     `json:"token"`
	UserID   uint64     `json:"userID"`
	Expires  *time.Time `json:"expires"`
	Verified bool       `json:"verified"`
}
