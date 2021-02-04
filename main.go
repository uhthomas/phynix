package main

import (
	"phynix/models"
	"phynix/realtime"
	"phynix/route"
	"raiki"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// m := gomail.NewMessage()
	// m.SetHeader("From", "noreply@phynix.io")
	// m.SetHeader("To", "m6f7262@gmail.com")
	// m.SetBody("text/plain", "Test from noreply.")

	// mail.Send(m)

	// models.DB.Create(&models.Playlist{
	// 	Name:   "Nice kek",
	// 	UserID: 1,
	// 	Active: true,
	// 	Items: []models.PlaylistItem{
	// 		{Title: "ok", Artist: "no", MediaID: 5, Order: 0},
	// 		{Title: "ok2", Artist: "no2", MediaID: 7, Order: 1},
	// 		{Title: "ok3", Artist: "no3", MediaID: 9, Order: 2},
	// 	},
	// })

	socketServer := raiki.NewServer()

	var communities []models.Community
	if err := models.DB.Find(&communities).Error; err != nil {
		panic(err)
	}

	for _, community := range communities {
		realtime.NewCommunity(community.ID, socketServer)
	}

	s := &route.Server{
		Domain:                      "phynix.io",
		CookieExpiration:            365 * 24 * time.Hour,
		CaptchaKey:                  "6LeW_hYTAAAAAC8-yhOI2b-ljh-s_TTBmtK8FNhm",
		UseCaptcha:                  true,
		EmailVerificationExpiration: 72 * time.Hour,
		Socket: socketServer,
	}

	s.Listen()
}
