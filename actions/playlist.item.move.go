package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func PlaylistItemMove(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID       uint64 `json:"id"`
		ItemID   uint64 `json:"itemID"`
		Position int    `json:"position"`
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

	var item models.PlaylistItem
	if tx.First(&item, data.ItemID).RecordNotFound() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemNonExist
	}

	var pErr error
	if err := playlist.ItemFunc(func(items []models.PlaylistItem) (payload []models.PlaylistItem) {
		if data.Position < 0 || data.Position > len(items)-1 {
			pErr = ErrPlaylistItemInvalidPosition
			return
		}

		for i, item := range items {
			if item.ID == data.ItemID {
				payload = append(items[:i], items[i+1:]...)
				payload = append(payload[:data.Position], append([]models.PlaylistItem{item}, payload[data.Position:]...)...)
				return
			}
		}
		pErr = ErrPlaylistItemNonExist
		return nil
	}); err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrMediaSave
	}

	if err := pErr; err != nil {
		return int(enum.ResponseCodeError), nil, err
	}

	tx.Commit()

	return int(enum.ResponseCodeOk), models.G{}, nil
}
