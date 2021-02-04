package actions

import (
	"encoding/json"
	"errors"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

var (
	PlaylistItemDeleteJSON     = errors.New("invalid json")
	PlaylistItemDeletePlaylist = errors.New("playlist doesn't exist")
	PlaylistItemDeleteItem     = errors.New("item doesn't exist")
	PlaylistItemDeleteOwner    = errors.New("you don't own this playlist")
	PlaylistItemDeleteLast     = errors.New("you cannot delete your last item")
	PlaylistItemDeleteSave     = errors.New("could not save")
)

func PlaylistItemDelete(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID     uint64 `json:"id"`
		ItemID uint64 `json:"itemID"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, PlaylistItemDeleteJSON
	}

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	tx := models.DB.Begin()

	var playlist models.Playlist
	if tx.First(&playlist, data.ID).RecordNotFound() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, PlaylistItemDeletePlaylist
	}

	if playlist.UserID != socket.ID() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, PlaylistItemDeleteOwner
	}

	var item models.PlaylistItem
	if tx.First(&item, data.ItemID).RecordNotFound() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, PlaylistItemDeleteItem
	}

	var pErr error
	if err := playlist.ItemFunc(func(items []models.PlaylistItem) (payload []models.PlaylistItem) {
		if len(items) == 1 {
			pErr = PlaylistItemDeleteLast
			return
		}

		for i, item := range items {
			if item.ID == data.ItemID {
				payload = append(items[:i], items[i+1:]...)
				return
			}
		}
		return nil
	}); err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, PlaylistItemDeleteSave
	}

	if pErr != nil {
		return int(enum.ResponseCodeBadRequest), nil, pErr
	}

	tx.Commit()

	return int(enum.ResponseCodeOk), models.G{}, nil
}
