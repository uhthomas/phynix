package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"phynix/enum"
	"phynix/models"
	"phynix/templates"
	"phynix/validation"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrSignupJson               = errors.New("could not decode json")
	ErrSignupQuery              = errors.New("unable to query database")
	ErrSignupCaptcha            = errors.New("could not verify captcha")
	ErrSignupInvalidUsername    = errors.New("invalid username")
	ErrSignupInvalidDisplayname = errors.New("invalid displayname")
	ErrSignupInvalidEmail       = errors.New("invalid email")
	ErrSignupInvalidPassword    = errors.New("invalid password")
	ErrSignupPassword           = errors.New("could not hash password")
	ErrSignupMail               = errors.New("could not send verification email")
)

func (s *Server) signup(res http.ResponseWriter, req *http.Request) {
	var data struct {
		Displayname string `json:"displayname"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Captcha     string `json:"captcha"`
	}

	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrSignupJson)
		return
	}

	if s.UseCaptcha {
		client := &http.Client{}
		resp, err := client.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
			"secret":   {s.CaptchaKey},
			"response": {data.Captcha},
			"remoteip": {strings.Split(req.RemoteAddr, ":")[0]},
		})
		if err != nil {
			s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupCaptcha)
			return
		}

		var out struct {
			Success bool     `json:"success"`
			Errors  []string `json:"error-codes"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupCaptcha)
			return
		}

		if !out.Success {
			s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupCaptcha)
			return
		}
	}

	var (
		displayname = strings.Trim(data.Displayname, " ")
		username    = strings.ToLower(data.Username)
		email       = strings.ToLower(data.Email)
		password    = data.Password
	)

	tx := models.DB.Begin()

	if !validation.Displayname(displayname) {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrSignupInvalidDisplayname)
		return
	}

	if !validation.Username(username) {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrSignupInvalidUsername)
		return
	}

	if !validation.Email(email) {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrSignupInvalidEmail)
		return
	}

	if !validation.Password(password) {
		s.writeResponse(res, enum.ResponseCodeBadRequest, nil, ErrSignupInvalidPassword)
		return
	}

	for k, v := range map[string]interface{}{
		"displayname": displayname,
		"username":    username,
		"email":       email,
	} {
		if err := tx.First(&models.User{}, k+" = ?", v).Error; err == nil {
			s.writeResponse(res, enum.ResponseCodeBadRequest, nil, fmt.Errorf("%s exists", k))
			return
		} else if err != gorm.ErrRecordNotFound {
			s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupQuery)
			return
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupPassword)
		return
	}

	u := models.User{
		Displayname: displayname,
		Username:    username,
		Email:       email,
		Hash:        hash,
		Tester:      true,
		Locale:      "en",
	}

	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupQuery)
		return
	}

	// vExp := time.Now().Add(s.EmailVerificationExpiration)

	// v := models.Verification{
	// 	Token:   models.Tokenize(string(u.ID)),
	// 	UserID:  u.ID,
	// 	Expires: &vExp,
	// }

	// if err := tx.Create(&v).Error; err != nil {
	// 	tx.Rollback()
	// 	s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupQuery)
	// 	return
	// }

	// m := gomail.NewMessage()
	// m.SetHeaders(map[string][]string{
	// 	"From":    {m.FormatAddress("support@phynix.io", "Phynix Support Bot")},
	// 	"To":      {m.FormatAddress(u.Email, u.Username)},
	// 	"Subject": {"Please verify your email"},
	// })

	// if err := mail.SendTemplate(m, "verification", models.G{
	// 	"token":    v.Token,
	// 	"username": u.Username,
	// }); err != nil {
	// 	tx.Rollback()
	// 	s.writeResponse(res, enum.ResponseCodeError, nil, ErrSignupMail)
	// 	return
	// }

	se := time.Now().Add(s.CookieExpiration)
	session := models.Session{
		Token:   models.Tokenize(string(u.ID)),
		UserID:  u.ID,
		Expires: &se,
	}

	if err := tx.Create(&session).Error; err != nil {
		tx.Rollback()
		res.Write([]byte(templates.Verify("Unable to generate session token", "")))
		return
	}

	tx.Commit()
	s.setCookie(res, "token", session.Token)
	s.writeResponse(res, enum.ResponseCodeOk, models.G{"token": session.Token}, nil)
}
