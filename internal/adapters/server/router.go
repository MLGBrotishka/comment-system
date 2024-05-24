package server

import (
	"net/http"
)

func NewRouter(router *http.ServeMux, handler http.Handler) {
	router.Handle("/query", handler)
}
