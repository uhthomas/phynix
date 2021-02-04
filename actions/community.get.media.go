package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"raiki"
	"time"
)

func CommunityGetMedia(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID uint64 `json:"id"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	db := models.DB

	var community models.Community
	if db.First(&community, data.ID).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityNonExist
	}

	c := realtime.GetCommunity(data.ID)
	if c == nil {
		return int(enum.ResponseCodeError), nil, ErrCommunityNonExist
	}

	media := c.Media
	if media != nil {
		media.Elapsed = time.Since(media.Timestamp).Seconds()
	}

	return int(enum.ResponseCodeOk), media, nil
}
