package actions

import (
	"encoding/json"
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"phynix/validation"
	"raiki"
	"strings"
)

func CommunityCreate(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
		Nsfw bool   `json:"nsfw"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	data.Name, data.Slug = strings.TrimSpace(data.Name), strings.TrimSpace(data.Slug)

	socket.Mu.Lock()
	defer socket.Mu.Unlock()

	tx := models.DB.Begin()
	u := models.User{
		Model: models.Model{
			ID: socket.ID(),
		},
	}

	if valid := validation.CommunityName(data.Name); !valid {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityName
	}

	if valid := validation.CommunitySlug(data.Slug); !valid {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunitySlug
	}

	if found := !tx.First(&models.Community{}, "slug = ?", data.Slug).RecordNotFound(); found {
		tx.Rollback()
		return int(enum.ResponseCodeBadRequest), nil, ErrCommunityExists
	}

	if count := tx.Model(&u).Association("Communities").Count(); count > 3 {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrCommunityLimit
	}

	community := models.Community{
		Name:    data.Name,
		Slug:    data.Slug,
		UserID:  socket.ID(),
		DjCycle: true,
		Nsfw:    data.Nsfw,
	}

	if err := tx.Model(&u).Association("Communities").Append(&community).Error; err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrCommunitySave
	}

	if err := tx.Model(&community).Association("Staff").Append(models.Staff{
		CommunityID: community.ID,
		UserID:      socket.ID(),
		Role:        enum.CommunityRoleHost,
	}).Error; err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrStaffSave
	}

	tx.Commit()

	realtime.NewCommunity(community.ID, socket.Server())

	return int(enum.ResponseCodeOk), community, nil
}
