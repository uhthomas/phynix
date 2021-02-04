package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func PlaylistGet(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID uint64 `json:"id"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	db := models.DB

	var playlist models.Playlist
	if db.First(&playlist, data.ID).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistNonExist
	}

	if playlist.UserID != socket.ID() {
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistOwner
	}

	return int(enum.ResponseCodeOk), playlist, nil
}
