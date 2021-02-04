package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func PlaylistDelete(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID uint64 `json:"id"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	tx := models.DB.Begin()
	u := models.User{
		Model: models.Model{
			ID: socket.ID(),
		},
	}

	var playlist models.Playlist
	if tx.First(&playlist, data.ID).RecordNotFound() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistNonExist
	}

	if playlist.UserID != socket.ID() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistOwner
	}

	if tx.Model(&u).Association("Playlists").Count() == 1 {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrPlaylistDeleteLast
	}

	if err := tx.Delete(&playlist).Error; err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrPlaylistDelete
	}

	tx.Commit()

	return int(enum.ResponseCodeOk), models.G{}, nil
}
