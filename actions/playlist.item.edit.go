package actions

import (
	"encoding/json"
	"errors"
	"phynix/enum"
	"phynix/models"
	"phynix/validation"
	"raiki"
	"strings"
)

var (
	ErrPlaylistItemEditJSON     = errors.New("invalid json")
	ErrPlaylistItemEditArtist   = errors.New("artist length invalid")
	ErrPlaylistItemEditTitle    = errors.New("title length invalid")
	ErrPlaylistItemEditPlaylist = errors.New("playlist doesn't exist")
	ErrPlaylistItemEditItem     = errors.New("item doesn't exist")
	ErrPlaylistItemEditOwner    = errors.New("you do not own this playlist")
	ErrPlaylistItemEditSave     = errors.New("could not save")
)

func PlaylistItemEdit(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID     uint64 `json:"id"`
		ItemID uint64 `json:"itemID"`
		Artist string `json:"artist"`
		Title  string `json:"title"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemEditJSON
	}

	data.Artist, data.Title = strings.TrimSpace(data.Artist), strings.TrimSpace(data.Title)

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	if !validation.PlaylistItemArtist(data.Artist) {
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemEditArtist
	}

	if !validation.PlaylistItemTitle(data.Title) {
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemEditTitle
	}

	tx := models.DB.Begin()

	var playlist models.Playlist
	if tx.First(&playlist, data.ID).RecordNotFound() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemEditPlaylist
	}

	if playlist.UserID != socket.ID() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemEditOwner
	}

	var item models.PlaylistItem
	if tx.First(&item, data.ItemID).RecordNotFound() {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemEditItem
	}

	item.Artist = data.Artist
	item.Title = data.Title

	if err := tx.Save(&item).Error; err != nil {
		return int(enum.ResponseCodeError), nil, ErrPlaylistItemEditSave
	}

	tx.Commit()

	return int(enum.ResponseCodeOk), item, nil
}
