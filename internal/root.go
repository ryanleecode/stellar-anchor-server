package internal

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

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

	apiRouter.
		HandleFunc("/authorizations", NewGetAuthHandler(authService)).
		Methods("GET")
	apiRouter.
		HandleFunc("/authorizations", NewPostAuthHandler(authService)).
		Methods("POST")

	return handlers.RecoveryHandler()(router)
}
