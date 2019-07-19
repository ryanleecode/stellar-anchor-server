package internal

import (
	"context"
	"crypto/rand"
	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/accounts"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/asset"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/authentication"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/random"
	stellarsdk "github.com/drdgvhbh/stellar-fi-anchor/internal/stellar-sdk"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"log"
	"net/http"
)

func Bootstrap(privateKey string, mnemonic string, db *gorm.DB) http.Handler {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatalln("cannot create accounts wallet")
	}
	passphrase := network.TestNetworkPassphrase
	fiKeyPair, err := keypair.Parse(privateKey)
	if err != nil {
		log.Fatalln("private key is not parsable")
	}

	ipcClient, err := rpc.DialIPC(context.TODO(), "/home/drd/.ethereum/goerli/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}
	ethIpcClient := ethclient.NewClient(ipcClient)
	ethIpcClient.ChainID(context.TODO()) // TODO

	db.AutoMigrate(asset.Asset{}, accounts.GenericAccount{})

	accountServices := []AccountService{accounts.NewService(wallet, db)}

	challengeTxFact := stellarsdk.NewChallengeTransactionFactory(
		passphrase,
		func() (s string, e error) {
			return random.NewGenerateString(random.NewGenerateBytes(rand.Read))(48)
		})
	authService := authentication.NewService(
		challengeTxFact, fiKeyPair.(*keypair.Full), passphrase)

	return NewRootHandler(authService, accountServices)
}
