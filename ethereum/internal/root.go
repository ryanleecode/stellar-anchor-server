package internal

import (
	"net/http"

	"github.com/drdgvhbh/stellar-anchor-server/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRootHandler(
	accountService AccountService,
) http.Handler {
	router := mux.NewRouter()

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(middleware.NewApplicationJSONMiddleware)

	apiRouter.
		HandleFunc("/deposit", NewGetDepositHandler(accountService)).
		Methods("GET")

	return handlers.RecoveryHandler()(router)
}
