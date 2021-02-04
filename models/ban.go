package models

import "time"

type Ban struct {
	Model
	BanneeID    uint64     `json:"banneeID"`
	BannerID    uint64     `json:"bannerID"`
	CommunityID uint64     `json:"communityID"`
	Reason      string     `json:"reason"`
	Expires     *time.Time `json:"expires"`
}
