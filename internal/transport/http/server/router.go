package server

// todo:  проблема, что чем больше роутов тем больше сюда будет стекать зависимостей, надо переделать чтобы сконфигурированный mux принимался параметром
// func NewRouter(authHandler *auth.Handler, mw middleware.MiddleWare) *http.ServeMux {
// 	mux := http.NewServeMux()
// 	auth.RegisterRoutes(mux, authHandler)

// 	return mux
// }
