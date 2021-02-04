package route

import (
	"net/http"
	"phynix/actions"
	"phynix/models"
	"phynix/templates"
)

func (s *Server) dashboard(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("token")
	if err != nil {
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	token := cookie.Value

	var sess models.Session
	if models.DB.First(&sess, "token = ?", token).RecordNotFound() {
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
		return
	}

	var u models.User
	if err := models.DB.First(&u, sess.UserID).Error; err != nil {
		res.Write([]byte("500"))
		return
	}

	_, communities, err := actions.CommunityList(nil, nil)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`500`))
		return
	}

	content := templates.Dashboard(u, communities.([]models.Community))
	res.Write([]byte(content))
}
