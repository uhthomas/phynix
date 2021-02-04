package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"raiki"
)

func CommunityGetUsers(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
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

	c := realtime.GetCommunity(data.ID)
	if c == nil {
		return int(enum.ResponseCodeError), nil, ErrCommunityNonExist
	}

	var users []models.User
	if err := db.Where(c.Users).Find(&users).Error; err != nil {
		return int(enum.ResponseCodeError), nil, nil
	}

	return int(enum.ResponseCodeOk), users, nil
}
