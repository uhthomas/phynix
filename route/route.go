package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"phynix/actions"
	"phynix/enum"
	"phynix/models"
	"phynix/realtime"
	"phynix/templates"
	"raiki"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Domain                      string
	CookieExpiration            time.Duration
	CaptchaKey                  string
	UseCaptcha                  bool
	EmailVerificationExpiration time.Duration
	Socket                      *raiki.Server
}

func (s *Server) Listen() {
	r := mux.NewRouter()

	s.Socket.SetAuthHandler(func(msg []byte) (uint64, interface{}, error) {
		var token string
		if err := json.Unmarshal(msg, &token); err != nil {
			return 0, nil, errors.New("could not decode token")
		}

		db := models.DB

		var sess models.Session
		if db.First(&sess, "token = ?", token).RecordNotFound() {
			return 0, nil, errors.New("invalid token")
		}

		var u models.User
		if err := db.
			Preload("Playlists").
			Preload("Playlists.Items").
			Preload("Communities").
			First(&u, sess.UserID).Error; err != nil {
			return 0, nil, errors.New("user not found")
		}

		if err := db.Where("user_id = ?", u.ID).Order("created_at desc").Limit(50).Find(&u.History).Error; err != nil {
			return 0, nil, errors.New("could not retrieve user history")
		}

		return u.ID, u, nil
	})

	s.Socket.SetCloseHandler(func(socket *raiki.ServerClient) {
		socket.Mu.Lock()
		defer socket.Mu.Unlock()
		u := realtime.NewUser(socket)
		if u == nil {
			fmt.Println("User doesn't exist!?")
			return
		}

		if c := u.Community; c != nil {
			c.Leave(socket.ID())
		} else {
			fmt.Println("Community doesn't exist!?")
		}
	})

	for n, f := range actions.SocketMap {
		s.Socket.Register(n, f)
	}

	r.HandleFunc("/", s.index)
	r.HandleFunc("/logout", s.logout)
	r.HandleFunc("/verify", s.verify)
	r.HandleFunc("/dashboard", s.dashboard)
	r.HandleFunc("/_/socket", s.Socket.ServeHTTP)
	r.HandleFunc("/_/signup", s.signup).Methods("POST")
	r.HandleFunc("/_/login", s.login).Methods("POST")
	r.HandleFunc("/_/bot/create", s.botCreate).Methods("POST")
	r.HandleFunc("/{community}", func(res http.ResponseWriter, req *http.Request) {
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

		if models.DB.First(&models.Community{}, "slug = ?", mux.Vars(req)["community"]).RecordNotFound() {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		content := templates.Loader()
		res.Write([]byte(content))
	})
	r.HandleFunc("/favicon.ico", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "/s/favicon.ico", http.StatusTemporaryRedirect)
	})
	r.PathPrefix("/s").Handler(http.StripPrefix("/s", http.FileServer(http.Dir("_/public"))))
	http.ListenAndServe(":9005", r)
}

func (s *Server) writeResponse(res http.ResponseWriter, status enum.ResponseCode, data interface{}, err error) {
	var e string
	if err != nil {
		e = err.Error()
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(status.Http())
	json.NewEncoder(res).Encode(struct {
		Status enum.ResponseCode `json:"status"`
		Data   interface{}       `json:"data,omitempty"`
		Error  string            `json:"error,omitempty"`
	}{status, data, e})
}

func (s *Server) setCookie(res http.ResponseWriter, name, value string) {
	http.SetCookie(res, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   fmt.Sprintf(".%s", s.Domain),
		Expires:  time.Now().Add(s.CookieExpiration),
		MaxAge:   int(s.CookieExpiration.Nanoseconds() / 1000000),
		Secure:   true,
		HttpOnly: true,
	})
}

func (s *Server) deleteCookie(res http.ResponseWriter, name string) {
	http.SetCookie(res, &http.Cookie{
		Name:     name,
		Value:    "deleted",
		Path:     "/",
		Domain:   fmt.Sprintf(".%s", s.Domain),
		Expires:  time.Now(),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	})
}

//
// func (s *Server) Signup(res http.ResponseWriter, req *http.Request) {
// 	var data struct {
// 		Username string `json:"username"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 		Captcha  string `json:"captcha"`
// 	}
//
// 	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
// 		s.writeResponse(res, enum.ResponseCodeError, nil, err)
// 		return
// 	}
// }
//
// func (s *Server) Login(res http.ResponseWriter, req *http.Request) {
// 	var data struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}
//
// 	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
// 		s.writeResponse(res, enum.ResponseCodeError, err, nil)
// 		return
// 	}
//
// 	user, err := db.GetUser(db.Q{"email": data.Email})
// 	if err == mgo.ErrNotFound {
// 		s.WriteResponse(res, enums.ResponseCodeForbidden, nil, errors.New("Wrong email or password"))
// 		return
// 	} else if err != nil {
// 		s.WriteResponse(res, enums.ResponseCodeError, nil, err)
// 		return
// 	}
//
// 	session, err := db.NewSession()
// }
