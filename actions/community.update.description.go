package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"phynix/validation"
	"raiki"
	"strings"
)

func CommunityUpdateDescription(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		ID          uint64 `json:"id"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	data.Description = strings.TrimSpace(data.Description)

	db := models.DB

	var community models.Community
	if db.First(&community, data.ID).RecordNotFound() {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityNonExist
	}

	if !community.HasPermission(socket.ID(), enum.CommunityRoleManager) {
		return int(enum.ResponseCodeForbidden), nil, ErrInsufficientPermission
	}

	if valid := validation.CommunityDescription(data.Description); !valid {
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityDescription
	}

	community.Description = data.Description

	if err := db.Save(&community).Error; err != nil {
		return int(enum.ResponseCodeError), nil, ErrCommunitySave
	}

	return int(enum.ResponseCodeOk), models.G{}, nil
}
