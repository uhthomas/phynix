package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func UserGet(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID uint64 `json:"id"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	db := models.DB

	var user models.User
	if db.First(&user, data.ID).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrUserNonExist
	}

	return int(enum.ResponseCodeOk), user, nil
}
