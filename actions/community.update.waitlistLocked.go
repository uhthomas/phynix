package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func CommunityUpdateWaitlistLocked(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID             uint64 `json:"id"`
		WaitlistLocked bool   `json:"WaitlistLocked"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	db := models.DB

	var community models.Community
	if db.First(&community, data.ID).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityNonExist
	}

	if !community.HasPermission(socket.ID(), enum.CommunityRoleManager) {
		return int(enum.ResponseCodeForbidden), nil, ErrInsufficientPermission
	}

	if community.WaitlistLocked == data.WaitlistLocked {
		return int(enum.ResponseCodeOk), models.G{}, nil
	}

	community.WaitlistLocked = data.WaitlistLocked

	if err := db.Save(&community).Error; err != nil {
		return int(enum.ResponseCodeError), nil, ErrCommunitySave
	}

	return int(enum.ResponseCodeOk), models.G{}, nil
}
