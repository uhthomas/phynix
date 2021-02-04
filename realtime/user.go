package realtime

import (
	"raiki"
	"sync"
)

type User struct {
	sync.Mutex
	Socket    *raiki.ServerClient
	Community *Community
}

var Users = map[uint64]*User{}

func NewUser(socket *raiki.ServerClient) *User {
	if u, ok := Users[socket.ID()]; ok {
		return u
	}

	u := &User{
		Socket:    socket,
		Community: nil,
	}

	Users[socket.ID()] = u
	return u
}

func (u *User) Hijack(socket *raiki.ServerClient) *User {
	u.Lock()
	defer u.Unlock()
	u.Socket.Close("Session replaced")
	u.Socket = socket
	return u
}
