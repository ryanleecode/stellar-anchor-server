package internal

import (
	"crypto/rand"
	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/drdgvhbh/stellar-fi-anchor/authentication/internal/random"
	"github.com/stellar/go/keypair"
	"log"
	"net/http"
)

type BootstrapParams interface {
	NetworkPassphrase() string
	Mnemonic() string
}

func Bootstrap(params BootstrapParams) http.Handler {
	wallet, err := hdwallet.NewFromMnemonic(params.Mnemonic())
	if err != nil {
		log.Fatalln("cannot create accounts wallet")
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	firstAccount, err := wallet.Derive(path, false)
	pk, err := wallet.PrivateKeyBytes(firstAccount)
	if err != nil {
		log.Fatalln("failed to derive private key bytes")
	}

	var pkBytes [32]byte
	copy(pkBytes[:], pk)
	fiKeyPair, err := keypair.FromRawSeed(pkBytes)
	if err != nil {
		log.Fatalln("failed to create keypair")
	}

	challengeTxFact := NewChallengeTransactionFactory(
		params.NetworkPassphrase(),
		func() (s string, e error) {
			return random.NewGenerateString(random.NewGenerateBytes(rand.Read))(48)
		})
	authService := NewService(
		challengeTxFact, fiKeyPair, params.NetworkPassphrase())

	return NewRootHandler(authService)
}