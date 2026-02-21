package server

import (
	"auth-test/internal/transport/http/auth"
	"net/http"
)

// todo:  проблема, что чем больше роутов тем больше сюда будет стекать зависимостей, надо переделать чтобы сконфигурированный mux принимался параметром
func NewRouter(authHandler *auth.Handler) http.Handler {
	mux := http.NewServeMux()
	auth.RegisterRoutes(mux, authHandler)
	return mux
}
