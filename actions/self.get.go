package actions

import (
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func SelfGet(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	db := models.DB

	var user models.User
	if db.Preload("Playlists").Preload("Playlists.Items").First(&user, socket.ID()).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrUserNonExist
	}

	return int(enum.ResponseCodeOk), user, nil
}
