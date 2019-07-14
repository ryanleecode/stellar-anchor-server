package internal

import (
	"crypto/rand"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"net/http"
	"stellar-fi-anchor/internal/authentication"
	"stellar-fi-anchor/internal/random"
	"stellar-fi-anchor/internal/stellar-sdk"

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
	fiKeyPair, err := keypair.Parse("SA4VF5RNXMFWS4JPLXDRP3D3SLSKMAZMCCXYC24LXMXUVYJLBN3F2ISY")
	if err != nil {
		panic(err)
	}

	client := horizonclient.DefaultTestNetClient
	clientWrpr := stellarsdk.NewClient(client)
	challengeTxFact := stellarsdk.NewChallengeTransactionFactory(
		network.TestNetworkPassphrase,
		func() (s string, e error) {
			return random.NewGenerateString(random.NewGenerateBytes(rand.Read))(48)
		})
	authService := authentication.NewService(
		clientWrpr, challengeTxFact, fiKeyPair.(*keypair.Full), network.TestNetworkPassphrase)

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
