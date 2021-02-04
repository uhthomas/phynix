package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func PlaylistActivate(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID uint64 `json:"id"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	tx := models.DB.Begin()

	var playlist models.Playlist
	if tx.First(&playlist, data.ID).RecordNotFound() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistNonExist
	}

	if playlist.UserID != socket.ID() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistOwner
	}

	if playlist.Active {
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistActive
	}

	if err := tx.Model(models.Playlist{}).Where("user_id = ? AND active = ?", socket.ID(), true).Updates(models.G{"active": false}).Error; err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrPlaylistSave
	}

	if err := tx.Model(&playlist).Update("active", true).Error; err != nil {
		return int(enum.ResponseCodeError), nil, ErrPlaylistSave
	}

	tx.Commit()

	return int(enum.ResponseCodeOk), models.G{}, nil
}
