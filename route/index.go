package route

import (
	"net/http"
	"phynix/actions"
	"phynix/models"
	"phynix/templates"
)

func (s *Server) index(res http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("token")
	if err == nil {
		token := cookie.Value

		var sess models.Session
		if !models.DB.First(&sess, "token = ?", token).RecordNotFound() {
			http.Redirect(res, req, "/dashboard", http.StatusTemporaryRedirect)
			return
		}
	}

	_, communities, err := actions.CommunityList(nil, nil)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`500`))
		return
	}

	content := templates.Index(communities.([]models.Community))
	res.Write([]byte(content))
}
