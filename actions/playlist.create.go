package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"phynix/validation"
	"raiki"
	"strings"
)

func PlaylistCreate(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	data.Name = strings.TrimSpace(data.Name)

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	tx := models.DB.Begin()
	u := models.User{
		Model: models.Model{
			ID: socket.ID(),
		},
	}

	if valid := validation.PlaylistName(data.Name); !valid {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistName
	}

	if tx.Model(&u).Association("Playlists").Count() >= 25 {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrPlaylistLimit
	}

	playlist := models.Playlist{
		Name:   data.Name,
		Active: true,
	}

	if err := tx.Model(models.Playlist{}).Where("user_id = ? AND active = ?", socket.ID(), true).Updates(models.G{"active": false}).Error; err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrPlaylistSave
	}

	if err := tx.Model(&u).Association("Playlists").Append(&playlist).Error; err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrPlaylistSave
	}

	tx.Commit()

	return int(enum.ResponseCodeOk), playlist, nil
}
