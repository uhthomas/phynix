package models

import "phynix/enum"

type User struct {
	Model
	Displayname string          `json:"displayname"`
	Username    string          `json:"username" gorm:"primary_key"`
	Email       string          `json:"-" gorm:"primary_key"`
	Hash        []byte          `json:"-"`
	GlobalRole  enum.GlobalRole `json:"globalRole"`
	Tester      bool            `json:"tester,omitempty"`
	Champion    bool            `json:"champion,omitempty"`
	Translator  bool            `json:"translator,omitempty"`
	Points      int             `json:"points"`
	FacebookId  *string         `json:"-" gorm:"primary_key"`
	TwitterId   *string         `json:"-" gorm:"primary_key"`
	Locale      string          `json:"locale"`

	Playlists   []Playlist  `json:"playlists,omitempty"`
	Communities []Community `json:"communities,omitempty"`
	History     []History   `json:"history,omitempty"`
	Chat        []Chat      `json:"chat,omitempty"`
	Sessions    []Session   `json:"sessions,omitempty"`
}
