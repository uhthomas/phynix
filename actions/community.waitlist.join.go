package actions

import (
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"raiki"
)

func CommunityWaitlistJoin(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	db := models.DB

	var playlist models.Playlist
	if db.First(&playlist, "user_id = ? and active = ?", socket.ID(), true).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityWaitlistJoinPlaylist
	}

	if db.Model(&playlist).Association("Items").Count() < 1 {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityWaitlistJoinPlaylistItems
	}

	u := realtime.NewUser(socket)

	if u == nil {
		return int(enum.ResponseCodeError), nil, ErrUserNonExist
	}

	c := u.Community

	if c == nil {
		return int(enum.ResponseCodeError), nil, ErrCommunityNonExist
	}

	return int(enum.ResponseCodeOk), models.G{
		"success": c.JoinWaitlist(socket.ID()),
	}, nil
}
