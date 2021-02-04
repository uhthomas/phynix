package models

import (
	"phynix/enum"
	"time"
)

type Staff struct {
	Model
	CommunityID uint64             `json:"communityID"`
	UserID      uint64             `json:"userID"`
	User        *User              `json:"user,omitempty"`
	Role        enum.CommunityRole `json:"role"`
	Expires     *time.Time         `json:"expires"`
}
