package route

import (
	"encoding/json"
	"errors"
	"net/http"
	"phynix/enum"
	"phynix/models"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLoginJson        = errors.New("could not decode json")
	ErrLoginQuery       = errors.New("unable to query database")
	ErrLoginInvalid     = errors.New("incorrect email or password")
	ErrLoginNotVerified = errors.New("account not verified")
	ErrLoginSession     = errors.New("unable to generate session token")
)

func (s *Server) login(res http.ResponseWriter, req *http.Request) {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrLoginJson)
		return
	}

	var user models.User
	if err := models.DB.First(&user, "email = ?", data.Email).Error; err == gorm.ErrRecordNotFound {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrLoginInvalid)
		return
	} else if err != nil {
		s.writeResponse(res, enum.ResponseCodeError, nil, ErrLoginQuery)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Hash, []byte(data.Password)); err != nil {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrLoginInvalid)
		return
	}

	if err := models.DB.First(&models.Verification{}, "user_id = ?", user.ID).Error; err == nil {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrLoginNotVerified)
		return
	} else if err != nil && err != gorm.ErrRecordNotFound {
		s.writeResponse(res, enum.ResponseCodeError, nil, ErrLoginQuery)
	}

	var sess models.Session
	if err := models.DB.First(&sess, "user_id = ?", user.ID).Error; err == gorm.ErrRecordNotFound {
		se := time.Now().Add(s.CookieExpiration)
		sess = models.Session{
			Token:   models.Tokenize(string(user.ID)),
			UserID:  user.ID,
			Expires: &se,
		}

		if err := models.DB.Create(sess).Error; err != nil {
			s.writeResponse(res, enum.ResponseCodeError, nil, ErrLoginSession)
			return
		}
	} else if err != nil {
		s.writeResponse(res, enum.ResponseCodeError, nil, ErrLoginQuery)
		return
	}

	s.setCookie(res, "token", sess.Token)
	s.writeResponse(res, enum.ResponseCodeOk, models.G{"token": sess.Token}, nil)
}
