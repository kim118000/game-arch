package handler

import "net/http"

func GetServerRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/v1/zone", HttpHeadDecorate(ZoneListHandler))
	router.Handle("/v1/login", HttpHeadDecorate(LoginHandler))
	return router
}
