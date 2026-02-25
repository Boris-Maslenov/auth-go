package auth

import (
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /sign-up", h.SignUp)
	mux.HandleFunc("POST /sign-in", h.SignIn)
	// mux.Handle("POST /sign-in", mw(http.HandlerFunc(h.SignIn)))
}
