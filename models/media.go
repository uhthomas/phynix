package models

type Media struct {
	Model
	Type      int    `json:"type"`
	ContentID string `json:"contentID"`
	Image     string `json:"image"`
	Duration  int    `json:"duration"`
	Artist    string `json:"artist"`
	Title     string `json:"title"`
	Blurb     string `json:"blurb"`
	Plays     int    `json:"plays"`
	Woots     int    `json:"woots"`
	Mehs      int    `json:"mehs"`
	Saves     int    `json:"saves"`
	Playlists int    `json:"playlists"`
}
