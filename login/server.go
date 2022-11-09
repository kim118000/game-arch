package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/kim118000/core/toolkit"
	"github.com/kim118000/core/toolkit/file"
	"github.com/kim118000/login/internal/conf"
	"github.com/kim118000/login/internal/handler"
	"github.com/kim118000/login/internal/service"
	"net/http"
)

func main() {

	confFile := file.PathJoin("conf/conf.toml")
	if _, err := toml.DecodeFile(confFile, conf.Config); err != nil {
		fmt.Println(err)
		return
	}
	service.InitService(conf.Config)

	serve := &http.Server{
		Handler: handler.GetServerRouter(),
		Addr:    fmt.Sprintf(":%d", conf.Config.ServerPort),
	}

	go serve.ListenAndServe()

	toolkit.RegisterSignal(func() {
		service.StopService()
	})
}
