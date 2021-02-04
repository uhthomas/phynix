package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"raiki"
)

func CommunityVoteSave(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID uint64 `json:"id"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	u := realtime.NewUser(socket)

	if u == nil {
		return int(enum.ResponseCodeError), nil, ErrUserNonExist
	}

	c := u.Community

	if c == nil {
		return int(enum.ResponseCodeError), nil, ErrCommunityNonExist
	}

	if c.Media == nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityMedia
	}

	payload, _ := json.Marshal(models.G{
		"id":        data.ID,
		"contentID": c.Media.Item.ContentID,
		"type":      c.Media.Item.Type,
		"bottom":    true,
		"artist":    c.Media.Item.Artist,
		"title":     c.Media.Item.Title,
	})

	socket.Mu.Unlock()
	s, d, err := PlaylistInsert(socket, payload)
	socket.Mu.Lock()
	if err != nil {
		return s, d, err
	}

	c.Vote(socket.ID(), enum.VoteTypeSave)

	return int(enum.ResponseCodeOk), d, nil
}
