package auth

import (
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /sign-up", h.SignUp)
	mux.HandleFunc("POST /sign-in", h.SignIn)
	mux.HandleFunc("GET /refresh", h.Refresh)
}
