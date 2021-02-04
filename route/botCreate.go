package route

import (
	"fmt"
	"net/http"
	"phynix/enum"
	"phynix/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) botCreate(res http.ResponseWriter, req *http.Request) {
	tx := models.DB.Begin()

	hash, err := bcrypt.GenerateFromPassword([]byte("bot_account"), 10)
	if err != nil {
		s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupPassword)
		return
	}

	a := fmt.Sprintf("bot%s", models.Tokenize("bot"))

	u := models.User{
		Displayname: a,
		Username:    a,
		Email:       fmt.Sprintf("%s@phynix.io", a),
		Hash:        hash,
		Tester:      true,
		Locale:      "en",
	}

	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupQuery)
		return
	}

	se := time.Now().Add(s.CookieExpiration)
	session := models.Session{
		Token:   models.Tokenize(string(u.ID)),
		UserID:  u.ID,
		Expires: &se,
	}

	if err := tx.Create(&session).Error; err != nil {
		tx.Rollback()
		s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupQuery)
		return
	}

	tx.Commit()

	s.setCookie(res, "token", session.Token)

	s.writeResponse(res, enum.ResponseCodeOk, models.G{
		"token": session.Token,
	}, nil)
}
