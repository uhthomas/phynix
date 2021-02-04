package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func CommunityGetHistory(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
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

	var history []models.History
	if err := db.Where("community_id = ?", community.ID).Order("updated_at desc").Limit(50).Find(&history).Error; err != nil {
		return int(enum.ResponseCodeError), nil, ErrHistoryGet
	}

	return int(enum.ResponseCodeOk), history, nil
}
