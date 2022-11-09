package handler

import "net/http"

func HttpHeadDecorate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")
		handler.ServeHTTP(w, r)
	})
}
