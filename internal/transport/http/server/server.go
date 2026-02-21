package server

import "net/http"

func NewServer(addr string, router http.Handler) *http.Server {
	return &http.Server{Addr: addr, Handler: router}
}
