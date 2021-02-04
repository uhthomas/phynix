package route

import (
	"net/http"
	"phynix/models"
	"phynix/templates"
	"time"

	"github.com/jinzhu/gorm"
)

func (s *Server) verify(res http.ResponseWriter, req *http.Request) {
	token := req.URL.Query()["token"]

	if len(token) == 0 {
		res.Write([]byte(templates.Verify("Invalid token", "")))
		return
	}

	tx := models.DB.Begin()

	var v models.Verification
	if err := tx.First(&v, "token = ?", token).Error; err == gorm.ErrRecordNotFound {
		res.Write([]byte(templates.Verify("Invalid token", "")))
		return
	} else if err != nil {
		res.Write([]byte(templates.Verify("Unable to query database", "")))
		return
	}

	if v.Verified {
		res.Write([]byte(templates.Verify("Already verified", "")))
		return
	}

	if v.Expires != nil && (*v.Expires).Before(time.Now()) {
		res.Write([]byte(templates.Verify("Token expired", "")))
		return
	}

	se := time.Now().Add(s.CookieExpiration)
	session := models.Session{
		Token:   models.Tokenize(string(v.UserID)),
		UserID:  v.UserID,
		Expires: &se,
	}

	if err := tx.Create(&session).Error; err != nil {
		tx.Rollback()
		res.Write([]byte(templates.Verify("Unable to generate session token", "")))
		return
	}

	if err := tx.Model(&models.Verification{}).Where("id = ?", v.ID).Updates("verified", true).Error; err != nil {
		tx.Rollback()
		res.Write([]byte(templates.Verify("Unable to invalidate verification session", "")))
		return
	}

	tx.Commit()
	s.setCookie(res, "token", session.Token)
	res.Write([]byte(templates.Verify(`Thanks for verifying your email!<br>You'll be redirected to the <a href="/dashboard">dashboard</a> in a few seconds`, session.Token)))
}
