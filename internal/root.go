package internal

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"net/http"
	"stellar-fi-anchor/internal/authorization"
	"stellar-fi-anchor/internal/stellar"

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
	var seed [32]byte
	copy(seed[:], []byte("SD5E2MXKN2MEOILRZ5TAXCNDC6ZSF2UZR7GDDMZCVRJQT3Q6BJCGYRTM"))
	fiKeyPair, err := keypair.FromRawSeed(seed)
	if err != nil {
		panic(err)
	}

	client := horizonclient.DefaultTestNetClient
	clientWrpr := stellar.NewClient(client)
	authService := authorization.NewService(clientWrpr, stellar.BuildChallengeTransaction, fiKeyPair)

	router := mux.NewRouter()
	router.Use(ContentType)
	router.Use(IDContext)
	router.Use(MethodContext)
	apiVersionRouter := router.PathPrefix("/v1").Subrouter()
	apiVersionRouter.Use(NewResponseWriter(newResponseWriter))

	apiVersionRouter.
		HandleFunc("/authorizations", NewGetAuthHandler(authService)).
		Methods("GET")
	apiVersionRouter.
		HandleFunc("/authorizations", NewPostAuthHandler(authService)).
		Methods("POST")

	return handlers.RecoveryHandler()(router)
}
