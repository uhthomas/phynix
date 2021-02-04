package actions

import (
	"encoding/json"
	"phynix/downloader"
	"phynix/enum"
	"phynix/models"
	"raiki"
	"strings"
)

func PlaylistInsert(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID        uint64         `json:"id"`
		ContentID string         `json:"contentID"`
		Type      enum.MediaType `json:"type"`
		Bottom    bool           `json:"bottom"`

		Artist *string `json:"artist"`
		Title  *string `json:"title"`
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

	if models.DB.Model(&playlist).Association("Items").Count() >= 500 {
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemLimit
	}

	if !models.DB.First(&models.PlaylistItem{}, "playlist_id = ? and content_id =? and type = ?", playlist.ID, data.ContentID, data.Type).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrPlaylistItemExists
	}

	var (
		media models.Media
		err   error
	)

	switch data.Type {
	case enum.MediaTypeYoutube:
		media, err = downloader.Youtube(data.ContentID)
	case enum.MediaTypeSoundcloud:
		media, err = downloader.Soundcloud(data.ContentID)
	default:
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrMediaType
	}

	if err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrMediaDownload
	}

	item := models.PlaylistItem{
		Type:       media.Type,
		ContentID:  media.ContentID,
		Image:      media.Image,
		Duration:   media.Duration,
		Artist:     media.Artist,
		Title:      media.Title,
		PlaylistID: playlist.ID,
		MediaID:    media.ID,
	}

	if data.Artist != nil {
		item.Artist = strings.TrimSpace(*data.Artist)
	}

	if data.Title != nil {
		item.Title = strings.TrimSpace(*data.Title)
	}

	// if err := tx.Create(&item).Error; err != nil {
	// 	tx.Rollback()
	// 	return int(enum.ResponseCodeError), nil, ErrMediaSave
	// }

	var i *models.PlaylistItem

	if err := playlist.ItemFunc(func(items []models.PlaylistItem) []models.PlaylistItem {
		var payload []models.PlaylistItem
		if data.Bottom {
			payload = append(items, item)
			i = &payload[len(payload)-1]
		} else {
			payload = append([]models.PlaylistItem{item}, items...)
			i = &payload[0]
		}
		return payload
	}); err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrMediaSave
	}

	tx.Commit()

	return int(enum.ResponseCodeOk), i, nil
}
