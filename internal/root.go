package internal

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"net/http"
	"stellar-fi-anchor/internal/authentication"
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
	authService := authentication.NewService(clientWrpr, stellar.BuildChallengeTransaction, fiKeyPair)

	router := mux.NewRouter()

	router.HandleFunc("/stellar.toml", func(w http.ResponseWriter, r *http.Request) {
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
