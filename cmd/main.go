// Package stellaranchor Stellar FI Anchor Emulator API.
//
// FI Anchor Emulator for the Stellar Network
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v1
//     Version: 0.0.2
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Ryan Lee<ryanleecode@gmail.com> http://drdgvhbh.io
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"crypto/rand"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"net/http"
	"os"
	"stellar-fi-anchor/internal/authentication"
	"stellar-fi-anchor/internal/random"
	stellarsdk "stellar-fi-anchor/internal/stellar-sdk"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"stellar-fi-anchor/internal"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalln("env variable PORT not defined")
	}
	privateKey, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		log.Fatalln("env variable PRIVATE_KEY not defined")
	}
	passphrase := network.TestNetworkPassphrase

	fiKeyPair, err := keypair.Parse(privateKey)
	if err != nil {
		log.Fatalln("private key is not parsable")
	}

	challengeTxFact := stellarsdk.NewChallengeTransactionFactory(
		passphrase,
		func() (s string, e error) {
			return random.NewGenerateString(random.NewGenerateBytes(rand.Read))(48)
		})
	authService := authentication.NewService(
		challengeTxFact, fiKeyPair.(*keypair.Full), passphrase)

	rootHandler := internal.NewRootHandler(authService)

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %d", 8000)
	log.Fatal(server.ListenAndServe())
}
