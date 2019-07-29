package internal

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRootHandler(
	authService AuthenticationService,
) http.Handler  {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.
		HandleFunc("/authentications", ApplicationJSON(NewGetAuthHandler(authService))).
		Methods("GET")
	apiRouter.
		HandleFunc("/authentications", UnsupportedMediaType(
			[]string{ "application/json", "application/x-www-form-urlencoded" },
		ApplicationJSON(
			PostAuthResponseWriter(
				NewPostAuthRequestValidator(
					NewPostAuthHandler(authService)))))).
		Methods("POST")

	return handlers.RecoveryHandler()(router)
}