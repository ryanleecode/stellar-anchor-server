package internal

import (
	"context"
	"crypto/rand"
	"fmt"
	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/accounts"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/asset"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/authentication"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/random"
	stellarsdk "github.com/drdgvhbh/stellar-fi-anchor/internal/stellar-sdk"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"log"
	"net/http"
)

func Bootstrap(privateKey string, mnemonic string, db *gorm.DB, rpcClient *rpc.Client) http.Handler {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatalln("cannot create accounts wallet")
	}
	passphrase := network.TestNetworkPassphrase
	fiKeyPair, err := keypair.Parse(privateKey)
	if err != nil {
		log.Fatalln("private key is not parsable")
	}

	ethClient := ethclient.NewClient(rpcClient)
	headers := make(chan *types.Header)
	sub, err := ethClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	go (func() {
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case header := <-headers:
				block, err := ethClient.BlockByHash(context.Background(), header.Hash())
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("--------------------------------------------------------------")
				fmt.Println(block.Hash().Hex())        // 0xbc10defa8dda384c96a17640d84de5578804945d347072e091b4e5f390ddea7f
				fmt.Println(block.Number().Uint64())   // 3477413
				fmt.Println(block.Time())              // 1529525947
				fmt.Println(block.Nonce())             // 130524141876765836
				fmt.Println(len(block.Transactions())) // 7
				for _, tx := range block.Transactions() {
					fmt.Println("TRANSACTION")
					fmt.Println(tx.Hash().Hex())
					fmt.Println(tx.Value())
					fmt.Println(tx.To().Hex())
				}
			}
		}
	})()

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
