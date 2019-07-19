package internal

import (
	"github.com/drdgvhbh/stellar-fi-anchor/internal/authentication"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/stellar/go/xdr"
	"net/http"
)

type AuthenticationService interface {
	BuildSignEncodeChallengeTransactionForAccount(id string) (string, error)
	ValidateClientSignedChallengeTransaction(
		txe *xdr.TransactionEnvelope) []error
	Authenticate(txe *xdr.TransactionEnvelope) (string, error)
	IsAuthorized(token string) bool
}

func NewRootHandler(
	authService AuthenticationService,
	accountServices []AccountService,
) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/.well-known/stellar.toml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "stellar.toml")
	})

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(ContentType)

	apiRouter.
		HandleFunc("/authorizations", authentication.NewGetAuthHandler(authService)).
		Methods("GET")
	apiRouter.
		HandleFunc("/authorizations", authentication.NewPostAuthHandler(authService)).
		Methods("POST")

	authMiddleware := NewBearerAuthMiddleware(authService)
	depositRouter := apiRouter.PathPrefix("/deposit").Subrouter()
	depositRouter.Use(authMiddleware)

	depositRouter.
		HandleFunc("", NewGetDepositHandler(accountServices)).
		Methods("GET")

	return handlers.RecoveryHandler()(router)
}
