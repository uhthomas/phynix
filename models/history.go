package models

import "time"

type History struct {
	Model
	CommunityID uint64    `json:"communityID"`
	UserID      uint64    `json:"userID"`
	User        User      `json:"user"`
	MediaID     uint64    `json:"mediaID"`
	Timestamp   time.Time `json:"timestamp"`
	Type        int       `json:"type"`
	ContentID   string    `json:"contentID"`
	Image       string    `json:"image"`
	Duration    int       `json:"duration"`
	Title       string    `json:"title"`
	Artist      string    `json:"artist"`
	Woots       int       `json:"woots"`
	Mehs        int       `json:"mehs"`
	Saves       int       `json:"saves"`
	Population  int       `json:"population"`
	Skipped     bool      `json:"skipped"`
}
