package route

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) logout(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Domain:   fmt.Sprintf(".%s", s.Domain),
		Expires:  time.Now(),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	})

	http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
}
