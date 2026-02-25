package middleware

import (
	"errors"
	"net/http"
	"strings"
)

type MiddleWare = func(next http.Handler) http.Handler

type Authenticator interface {
	ParseToken(token string) (int64, error)
}

func Auth(a Authenticator) MiddleWare {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := getTokenFromHeader(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			_, err = a.ParseToken(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getTokenFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	tokenParts := strings.Split(header, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" || len(tokenParts[1]) == 0 {
		return "", errors.New("invalid token")
	}

	return tokenParts[1], nil
}
