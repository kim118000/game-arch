package handler

import (
	"encoding/json"
	"github.com/kim118000/core/pkg/logger"
	"github.com/kim118000/login/internal/data"
	"net/http"
)

var ZoneListHandler = new(ZoneList)
type ZoneList struct {}

func (z *ZoneList) ServeHTTP(w http.ResponseWriter, r *http.Request) () {
	zone := []*data.ZoneData{
		{
			Id:   1,
			Name: "cn",
			Url:  "127.0.0.1:8080",
		},
		{
			Id:   2,
			Name: "en",
			Url:  "127.0.0.1:8080",
		},
	}

	data, _ := json.Marshal(zone)
	_, err := w.Write(data)
	if err != nil {
		logger.Log.Errorf("get zone list %s", err)
	}
}
