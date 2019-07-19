package internal

import (
	"net/http"
	"regexp"
	"strings"
)

var apiCallPrefix = regexp.MustCompile("/api/v[0-9]+/")

// ContentType injects application/json as the content type
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

type AuthService interface {
	IsAuthorized(token string) bool
}

func NewBearerAuthMiddleware(authService AuthService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			authHeader := request.Header.Get("authorization")
			token := strings.TrimPrefix(authHeader, "Bearer ")

			if authService.IsAuthorized(token) {
				next.ServeHTTP(writer, request)
			} else {
				writer.WriteHeader(http.StatusUnauthorized)
			}
		})
	}
}
