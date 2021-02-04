package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"raiki"
)

func CommunityGetStaff(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
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

	var staff []models.Staff
	if err := db.Where("community_id = ?", community.ID).Order("updated_at desc").Find(&staff).Error; err != nil {
		return int(enum.ResponseCodeError), nil, ErrStaffGet
	}

	return int(enum.ResponseCodeOk), staff, nil
}
