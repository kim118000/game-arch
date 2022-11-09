package handler

import (
	"encoding/json"
	"github.com/kim118000/login/internal/service"
	"net/http"
)

var LoginHandler = new(Login)

type Login struct{}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) () {
	username := r.URL.Query().Get("u")
	password := r.URL.Query().Get("p")
	service.Log.Infof("login user=%s pwd=%s", username, password)

	loginInfo := struct {
		Name string `json:"name"`
		Pwd  string `json:"pwd"`
	}{
		Name: username,
		Pwd:  password,
	}

	data, _ := json.Marshal(loginInfo)
	_, err := w.Write(data)
	if err != nil {
		service.Log.Errorf("login %s", err)
	}
}
