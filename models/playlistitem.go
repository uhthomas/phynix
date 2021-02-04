package models

type PlaylistItem struct {
	Model
	Type       int    `json:"type"`
	ContentID  string `json:"contentID"`
	Image      string `json:"image"`
	Duration   int    `json:"duration"`
	Artist     string `json:"artist"`
	Title      string `json:"title"`
	PlaylistID uint64 `json:"playlistID"`
	MediaID    uint64 `json:"mediaID"`
	Position   int    `json:"position"`
}
