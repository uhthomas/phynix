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

func ChatSend(socket *raiki.ServerClient, msg []byte) (int, interface{}, error) {
	var data struct {
		Emote   bool   `json:"emote"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(msg, &data); err != nil {
		return int(enum.ResponseCodeBadRequest), nil, ErrInvalidJson
	}

	data.Message = strings.TrimSpace(data.Message)

	u := realtime.NewUser(socket)

	if u == nil {
		return int(enum.ResponseCodeError), nil, ErrUserNonExist
	}

	c := u.Community

	if c == nil {
		return int(enum.ResponseCodeError), nil, ErrCommunityNonExist
	}

	tx := models.DB.Begin()

	var mute models.Mute
	if !tx.First(&models.Mute{}, "community_id = ? AND mutee_id = ?", c.ID, socket.ID()).RecordNotFound() {
		if mute.Expires == nil || mute.Expires.After(time.Now()) {
			tx.Rollback()
			return int(enum.ResponseCodeForbidden), nil, ErrMuted
		} else if err := tx.Delete(&mute).Error; err != nil {
			tx.Rollback()
			return int(enum.ResponseCodeError), nil, ErrMuteDelete
		}
	}

	message := data.Message
	if message == "" {
		return int(enum.ResponseCodeBadRequest), nil, ErrChatShort
	}

	if len(message) > 255 {
		message = message[:255]
	}

	chat := models.Chat{
		UserID:      socket.ID(),
		CommunityID: c.ID,
		Emote:       data.Emote,
		Message:     message,
	}

	if err := tx.Save(&chat).Error; err != nil {
		tx.Rollback()
		return int(enum.ResponseCodeError), nil, ErrChatSave
	}

	tx.Commit()

	socket.Emit(c.Room, raiki.Event{"chat.receive", chat})

	return int(enum.ResponseCodeOk), models.G{}, nil
}
