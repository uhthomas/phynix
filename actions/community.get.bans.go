package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func CommunityGetBans(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID uint64 `json:"id"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	db := models.DB

	var community models.Community
	if db.First(&community, data.ID).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityNonExist
	}

	var bans []models.Ban
	if err := db.Model(&community).Association("Bans").Find(&bans).Error; err != nil {
		return int(enum.ResponseCodeError), nil, ErrBanGet
	}

	return int(enum.ResponseCodeOk), bans, nil
}
