package internal

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var apiCallPrefix = regexp.MustCompile("/api/v[0-9]+/")

// ContentType injects application/json as the content type
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

// MethodContext injects the http route, along with its http method in the response
func MethodContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		ctx := context.WithValue(
			request.Context(),
			"method",
			fmt.Sprintf("%s.%s",
				strings.Replace(
					apiCallPrefix.ReplaceAllString(
						request.RequestURI, ""), "/", ".", -1),
				strings.ToLower(request.Method),
			)[4:])

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// IDContext injects a UUID as the request ID in the response
func IDContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.WithValue(request.Context(), "id", uuid.New().String())

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

// NewResponseWriter creates a middleware that intercepts another http.ResponseWriter
func NewResponseWriter(newResponseWriter func(http.ResponseWriter, *http.Request) http.ResponseWriter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(writer http.ResponseWriter, request *http.Request) {
				interceptedWriter := newResponseWriter(writer, request)

				next.ServeHTTP(interceptedWriter, request)
			})
	}
}
