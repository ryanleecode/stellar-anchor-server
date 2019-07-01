package internal

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	responseProperties = &Properties{
		APIVersion: "0.0.2",
	}
)

func newResponseWriter(
	w http.ResponseWriter,
	r *http.Request,
) http.ResponseWriter {
	return &ResponseWriter{
		ResponseProperties: *responseProperties,
		Writer:             w,
		Request:            r,
	}
}

func NewRootHandler() http.Handler {
	router := mux.NewRouter()
	router.Use(ContentType)
	router.Use(IDContext)
	router.Use(MethodContext)
	apiVersionRouter := router.PathPrefix("/v1").Subrouter()
	apiVersionRouter.Use(NewResponseWriter(newResponseWriter))

	return handlers.RecoveryHandler()(router)
}
