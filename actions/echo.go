package actions

import (
	"phynix/enum"
	"raiki"
)

func Echo(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	return int(enum.ResponseCodeOk), msg, nil
}
