package models

type Chat struct {
	Model
	UserID      uint64  `json:"userID"`
	User        *User   `json:"user,omitempty"`
	CommunityID uint64  `json:"communityID"`
	DeleterID   *uint64 `json:"deleterID,omitempty"`
	Emote       bool    `json:"emote"`
	Message     string  `json:"message"`
	Edited      bool    `json:"edited"`
}
