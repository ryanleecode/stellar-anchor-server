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
	"context"
	"crypto/rand"
	"fmt"
	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"net/http"
	"os"
	"stellar-fi-anchor/internal/accounts"
	"stellar-fi-anchor/internal/asset"
	"stellar-fi-anchor/internal/authentication"
	"stellar-fi-anchor/internal/random"
	stellarsdk "stellar-fi-anchor/internal/stellar-sdk"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
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
	ethMnemonic, ok := os.LookupEnv("ETHEREUM_MNEMONIC")
	if !ok {
		log.Fatalln("env variable ETHEREUM_MNEMONIC not defined")
	}
	wallet, err := hdwallet.NewFromMnemonic(ethMnemonic)
	if err != nil {
		log.Fatalln("cannot create accounts wallet")
	}
	print(wallet)

	db, err := gorm.Open(
		"postgres", "host=localhost port=6666 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err, "failed to open database")
	}
	defer func() {
		_ = db.Close()
	}()

	ipcClient, err := rpc.DialIPC(context.TODO(), "/home/drd/.ethereum/goerli/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}
	ethIpcClient := ethclient.NewClient(ipcClient)
	chainId, err := ethIpcClient.ChainID(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	println(chainId.String())

	db.AutoMigrate(asset.Asset{}, accounts.GenericAccount{})

	accountServices := []internal.AccountService{accounts.NewService(wallet, db)}

	challengeTxFact := stellarsdk.NewChallengeTransactionFactory(
		passphrase,
		func() (s string, e error) {
			return random.NewGenerateString(random.NewGenerateBytes(rand.Read))(48)
		})
	authService := authentication.NewService(
		challengeTxFact, fiKeyPair.(*keypair.Full), passphrase)

	rootHandler := internal.NewRootHandler(authService, accountServices)

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %d", 8000)
	log.Fatal(server.ListenAndServe())
}
