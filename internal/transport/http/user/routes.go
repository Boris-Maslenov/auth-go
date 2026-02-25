package user

import (
	"auth-test/internal/transport/http/middleware"
	"net/http"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[user1, user2, user3, user4]"))
}

func RegisterUserRoutes(mux *http.ServeMux, mw middleware.MiddleWare) {
	mux.Handle("GET /users", mw(http.HandlerFunc(getUsers)))
}
