package test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stellar/go/network"

	"github.com/drdgvhbh/stellar-fi-anchor/sdk"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal"

	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/suite"
)

type DepositEthereumSuite struct {
	suite.Suite
	clientKP  *keypair.Full
	server    *httptest.Server
	db        *gorm.DB
	rpcClient *rpc.Client
	apiClient *sdk.APIClient
}

func (s *DepositEthereumSuite) SetupSuite() {
	err := godotenv.Load("../.env.test")
	s.NoError(err)

	env := NewEnvironment()
	db, err := gorm.Open(
		"postgres", fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			env.DBHost(),
			env.DBPort(),
			env.DBUser(),
			env.DBName(),
			env.DBSSLMode(),
			env.DBPassword()))
	s.NoError(err)
	s.db = db

	rpcClient, err := rpc.DialHTTP(env.ethRPCEndpoint)
	s.NoError(err)
	s.rpcClient = rpcClient

	issuerKP, err := keypair.Random()
	s.NoError(err)
	client := horizonclient.DefaultTestNetClient
	_, err = client.Fund(issuerKP.Address())
	s.NoError(err)
}

func (s *DepositEthereumSuite) SetupTest() {
	mnemonic, err := hdwallet.NewMnemonic(128)
	s.NoError(err)

	logger := logrus.New()

	rootHandler := internal.Bootstrap(network.TestNetworkPassphrase, mnemonic, s.db, s.rpcClient, logger)
	s.server = httptest.NewServer(rootHandler)

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

	d, _, err := s.apiClient.AccountApi.Deposit(ctx, account, "ETH")
	s.NoError(err)

	s.Regexp(regexp.MustCompile("^0x[a-fA-F0-9]{40}$"), d.How)
}

func (s *DepositEthereumSuite) TestMultipleDepositCallsForSameAssetShouldReturnSameAddress() {
	ctx := context.Background()
	account := s.clientKP.Address()

	d, _, err := s.apiClient.AccountApi.Deposit(ctx, account, "ETH")
	s.NoError(err)

	firstAddress := d.How

	d, _, err = s.apiClient.AccountApi.Deposit(ctx, account, "ETH")
	s.NoError(err)

	secondAddress := d.How

	s.EqualValues(firstAddress, secondAddress)
}

func (s *DepositEthereumSuite) AfterTest(_, _ string) {
	s.server.Close()
}

func TestDepositEthereumSuite(t *testing.T) {
	suite.Run(t, new(DepositEthereumSuite))
}
