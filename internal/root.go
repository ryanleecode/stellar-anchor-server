package internal

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
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

func NewRootHandler(
	authService AuthenticationService,
) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/.well-known/stellar.toml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "stellar.toml")
	})

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(ContentType)
	apiRouter.Use(IDContext)
	apiRouter.Use(MethodContext)
	apiRouter.Use(NewResponseWriter(newResponseWriter))

	apiRouter.
		HandleFunc("/authorizations", NewGetAuthHandler(authService)).
		Methods("GET")
	apiRouter.
		HandleFunc("/authorizations", NewPostAuthHandler(authService)).
		Methods("POST")

	return handlers.RecoveryHandler()(router)
}
