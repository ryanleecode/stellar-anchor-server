package internal

import (
	"net/http"
)

func ApplicationJSON(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("content-type", "application/json")
		next.ServeHTTP(writer, request)
	})
}

func UnsupportedMediaType(allowedMediaTypes []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("content-type")
		for _, allowedType := range allowedMediaTypes {
			if contentType == allowedType {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.WriteHeader(http.StatusUnsupportedMediaType)
	})
}
