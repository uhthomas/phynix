package route

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"text/template"
)

func (s *Server) template(name string) func(res http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile(filepath.Join("_", "template", "web", name+".tmpl"))
	if err != nil {
		panic(err)
	}

	return func(res http.ResponseWriter, req *http.Request) {
		t, err := template.New(name).Parse(fmt.Sprintf("%s", content))
		if err != nil {
			res.Write([]byte(err.Error()))
			return
		}

		res.Header().Set("Content-Type", "text/html; encoding=utf-8")
		if err := t.Execute(res, nil); err != nil {
			res.Write([]byte(err.Error()))
			return
		}
	}
}
