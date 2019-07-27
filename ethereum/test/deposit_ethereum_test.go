package test

import (
	"context"
	"fmt"
	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/drdgvhbh/stellar-fi-anchor/authorization/internal"
	"github.com/drdgvhbh/stellar-fi-anchor/authorization/sdk"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
)

type DepositEthereumSuite struct {
	suite.Suite
	serverKP  *keypair.Full
	clientKP  *keypair.Full
	server    *httptest.Server
	db        *gorm.DB
	apiClient *sdk.APIClient
}

func (s *DepositEthereumSuite) SetupSuite() {
	err := godotenv.Load("../.env")
	s.NoError(err)
	serverKP, err := keypair.Random()
	s.NoError(err)
	client := horizonclient.DefaultTestNetClient
	_, err = client.Fund(serverKP.Address())
	s.NoError(err)
	s.serverKP = serverKP
}

func (s *DepositEthereumSuite) SetupTest() {
	mnemonic, err := hdwallet.NewMnemonic(128)
	s.NoError(err)
	serverPK := s.serverKP.Seed()
	db, err := gorm.Open(
		"postgres", "host=localhost port=6666 user=postgres dbname=postgres sslmode=disable")
	s.NoError(err)

	rpcClient, err := rpc.DialHTTP(os.Getenv("INFURA_URL"))
	s.NoError(err)
	rootHandler := internal.Bootstrap(serverPK, mnemonic, db, rpcClient)
	s.server = httptest.NewServer(rootHandler)
	s.db = db

	clientKP, err := keypair.Random()
	s.NoError(err)
	s.clientKP = clientKP

	apiClient := sdk.NewAPIClient(&sdk.Configuration{
		BasePath: fmt.Sprintf("%s/api", s.server.URL),
	})
	s.apiClient = apiClient
}

func (s *DepositEthereumSuite) TestDepositRouteWorks() {
	ctx := context.Background()
	account := s.clientKP.Address()
	tx, _, err := s.apiClient.AuthorizationApi.RequestAChallenge(ctx, account)
	s.NoError(err)

	var txe xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(tx.Transaction, &txe)
	s.NoError(err)

	b := &build.TransactionEnvelopeBuilder{E: &txe}
	b.Init()
	err = b.MutateTX(build.TestNetwork)
	s.NoError(err)
	err = b.Mutate(build.Sign{Seed: s.clientKP.Seed()})
	s.NoError(err)

	signedTxe, err := b.Base64()

	authToken, _, err := s.apiClient.
		AuthorizationApi.
		Authenticate(ctx, account, sdk.ChallengeTransaction{Transaction: signedTxe})
	s.NoError(err)

	authCtx := context.WithValue(ctx, sdk.ContextAccessToken, authToken.Token)
	d, _, err := s.apiClient.AccountApi.Deposit(authCtx, account, "ETH")
	s.NoError(err)

	s.Regexp(regexp.MustCompile("^0x[a-fA-F0-9]{40}$"), d.How)
}

func (s *DepositEthereumSuite) TestMultipleDepositCallsForSameAssetShouldReturnSameAddress() {
	ctx := context.Background()
	account := s.clientKP.Address()
	tx, _, err := s.apiClient.AuthorizationApi.RequestAChallenge(ctx, account)
	s.NoError(err)

	var txe xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(tx.Transaction, &txe)
	s.NoError(err)

	b := &build.TransactionEnvelopeBuilder{E: &txe}
	b.Init()
	err = b.MutateTX(build.TestNetwork)
	s.NoError(err)
	err = b.Mutate(build.Sign{Seed: s.clientKP.Seed()})
	s.NoError(err)

	signedTxe, err := b.Base64()

	authToken, _, err := s.apiClient.
		AuthorizationApi.
		Authenticate(ctx, account, sdk.ChallengeTransaction{Transaction: signedTxe})
	s.NoError(err)

	authCtx := context.WithValue(ctx, sdk.ContextAccessToken, authToken.Token)
	d, _, err := s.apiClient.AccountApi.Deposit(authCtx, account, "ETH")
	s.NoError(err)

	firstAddress := d.How

	d, _, err = s.apiClient.AccountApi.Deposit(authCtx, account, "ETH")
	s.NoError(err)

	secondAddress := d.How

	s.EqualValues(firstAddress, secondAddress)
}

func (s *DepositEthereumSuite) AfterTest(_, _ string) {
	s.server.Close()
	err := s.db.Close()
	s.NoError(err)
}

func TestDepositEthereumSuite(t *testing.T) {
	suite.Run(t, new(DepositEthereumSuite))
}
