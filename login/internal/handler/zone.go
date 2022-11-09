package handler

import (
	"encoding/json"
	"github.com/kim118000/login/internal/data"
	"github.com/kim118000/login/internal/service"
	"net/http"
)

var ZoneListHandler = new(Zone)
type Zone struct {
}

func (l *Zone) ServeHTTP(w http.ResponseWriter, r *http.Request) () {
	zone := []*data.ZoneData{
		{
			Id:   1,
			Name: "cn",
			Url:  "127.0.0.1:8999",
		},
		{
			Id:   2,
			Name: "en",
			Url:  "127.0.0.1:8999",
		},
	}

	data, _ := json.Marshal(zone)
	_, err := w.Write(data)
	if err != nil {
		service.Log.Errorf("get zone list %s", err)
	}
}
