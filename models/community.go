package models

import (
	"phynix/enum"
)

type Community struct {
	Model
	Name           string `json:"name"`
	Slug           string `json:"slug" gorm:"primary_key"`
	User           User   `json:"host"`
	UserID         uint64 `json:"-"`
	Description    string `json:"description"`
	WelcomeMessage string `json:"welcomeMessage"`
	WaitlistLocked bool   `json:"waitlistLocked"`
	DjCycle        bool   `json:"djCycle"`
	Nsfw           bool   `json:"nsfw"`

	Population int           `json:"population,omitempty" sql:"-"`
	Waitlist   int           `json:"waitlist,omitempty" sql:"-"`
	Media      PlaylistItem `json:"media,omitempty" sql:"-"`

	Staff   []Staff   `json:"staff,omitempty"`
	History []History `json:"history,omitempty"`
	Chat    []Chat    `json:"chat,omitempty"`
	Bans    []Ban     `json:"ban,omitempty"`
	Mutes   []Mute    `json:"mutes,omitempty"`
}

func (c Community) HasPermission(userID uint64, required enum.CommunityRole) bool {
	var user User
	if DB.First(&user, userID).RecordNotFound() {
		return false
	}

	if user.GlobalRole >= enum.GlobalRoleModerator {
		return true
	}

	if DB.First(&Staff{}, "community_id = ? and user_id = ? and role >= ?", c.ID, userID, required).RecordNotFound() {
		return false
	}

	return true
}

// func (c Community) AfterUpdate() error {
// 	go UploadToAlgolia("community", c)
// 	return nil
// }

// func (c Community) AfterSave() error {
// 	go UploadToAlgolia("community", c)
// 	return nil
// }
