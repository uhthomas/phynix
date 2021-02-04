package models

import "time"

type GlobalBan struct {
	Model
	BanneeID uint64     `json:"banneeID"`
	BannerID uint64     `json:"bannerID"`
	Reason   string     `json:"reason"`
	Expires  *time.Time `json:"expires"`
}
