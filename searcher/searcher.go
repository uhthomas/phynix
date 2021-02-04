package searcher

type Result struct {
	Type      int    `json:"type"`
	ContentID string `json:"contentID"`
	Image     string `json:"image"`
	Duration  int    `json:"duration"`
	Artist    string `json:"artist"`
	Title     string `json:"title"`
}
