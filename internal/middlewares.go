package internal

import (
	"net/http"
	"regexp"
)

var apiCallPrefix = regexp.MustCompile("/api/v[0-9]+/")

// ContentType injects application/json as the content type
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}
