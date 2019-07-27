package accounts

import (
	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/mock"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EthereumAccountServiceSuite struct {
	suite.Suite
	accountService *EthereumAccountService
	wallet         *mock.WalletMock
}

func (s *EthereumAccountServiceSuite) SetupSuite() {
	mocket.Catcher.Register()
}

func (s *EthereumAccountServiceSuite) SetupTest() {
	mocket.Catcher.Reset()
	mocket.Catcher.Logging = true

	db, err := gorm.Open(mocket.DriverName, "connection_string")
	s.NoError(err)
	s.wallet = &mock.WalletMock{}
	s.accountService = NewService(s.wallet, db)
}

/*func (s *EthereumAccountServiceSuite) TestNewAccount_AccountNumberIsIncrementedBasedOnAssetType() {
	lastAccountNumber := 10
	queryResponse := []map[string]interface{}{{
		"asset_type":         asset.Ethereum,
		"stellar_account_id": "GDCANY73N6JU5TGXQJI6NSDYKCPPNO5EXZNO4WOCKTQLM7MCMEQKQVFA",
		"number":             lastAccountNumber,
	}}
	mocket.Catcher.NewMock().
		WithQuery(`SELECT * FROM "accounts"  WHERE "accounts"."deleted_at" IS NULL AND (("accounts"."asset_type" = ETH)) ORDER BY number desc,"accounts"."id" ASC LIMIT 1`).
		WithReply(queryResponse)

	ethAddress := "0x87fa4adb38EF0bF29F3EE8035db10fd13B70c0d1"
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", lastAccountNumber+1))
	s.wallet.On("Derive", path, false).Return(
		accounts.Account{Address: common.HexToAddress(ethAddress)}, nil)

	createdAccount, err := s.accountService.GetDepositingAccount(
		"GAJN3BQGETOKZYI6BYNQJZ5ZQWURZX5TXNOEF7UPQAB4BHBCW5JEIOAD")
	s.NoError(err)

	assert.EqualValues(s.T(), lastAccountNumber+1, createdAccount.DepositInstructions())
}*/

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(EthereumAccountServiceSuite))
}
