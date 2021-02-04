package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"raiki"
	"strings"
	"time"
)

func CommunityJoin(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		Slug string `json:"slug"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	data.Slug = strings.TrimSpace(data.Slug)

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	db := models.DB

	var community models.Community
	if db.First(&community, "slug = ?", data.Slug).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityNonExist
	}

	var ban models.Ban
	if !db.First(&models.Ban{}, "community_id = ? and bannee_id = ?", community.ID, socket.ID()).RecordNotFound() {
		if ban.Expires == nil || ban.Expires.After(time.Now()) {
			return int(enum.ResponseCodeForbidden), models.G{"expires": ban.Expires}, ErrBanned
		} else if err := db.Delete(&ban).Error; err != nil {
			return int(enum.ResponseCodeError), nil, ErrBanDelete
		}
	}

	u := realtime.NewUser(socket)

	if u == nil {
		return int(enum.ResponseCodeError), nil, ErrUserNonExist
	}

	c := realtime.GetCommunity(community.ID)

	if c == nil {
		return int(enum.ResponseCodeError), nil, ErrCommunityNonExist
	}

	if c := u.Community; c != nil {
		c.Leave(socket.ID())
	}

	u.Community = c
	c.Join(socket.ID())
	socket.Join(c.Room)

	state, err := c.State()
	if err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityState
	}

	return int(enum.ResponseCodeOk), state, nil
}
