package models

import "time"

type Mute struct {
	Model
	MuteeID     uint64     `json:"mutteeID"`
	MuterID     uint64     `json:"muterID"`
	CommunityID uint64     `json:"communityID"`
	Reason      string     `json:"reason"`
	Expires     *time.Time `json:"expires"`
}
