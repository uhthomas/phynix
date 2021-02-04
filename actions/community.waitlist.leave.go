package actions

import (
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"raiki"
)

func CommunityWaitlistLeave(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
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

	return int(enum.ResponseCodeOk), models.G{
		"success": c.LeaveWaitlist(socket.ID()),
	}, nil
}
