package accounts

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/suite"
	"stellar-fi-anchor/internal/models"
	"stellar-fi-anchor/mock"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	accountService *Service
}

func (s *ServiceSuite) SetupSuite() {
	mocket.Catcher.Register()
}

func (s *ServiceSuite) SetupTest() {
	mocket.Catcher.Reset()
	mocket.Catcher.Logging = true

	db, err := gorm.Open(mocket.DriverName, "connection_string")
	s.NoError(err)
	walletMock := mock.WalletMock{}
	s.accountService = NewService(walletMock, db)
}

func (s *ServiceSuite) TestNewAccount_AccountNumberIsIncrementedBasedOnAssetType() {
	lastAccountNumber := 10
	response := []map[string]interface{}{{
		"asset_type":         models.Ethereum,
		"stellar_account_id": "GDCANY73N6JU5TGXQJI6NSDYKCPPNO5EXZNO4WOCKTQLM7MCMEQKQVFA",
		"number":             lastAccountNumber,
	}}
	mocket.Catcher.NewMock().
		WithQuery(`SELECT * FROM "accounts"  WHERE "accounts"."deleted_at" IS NULL AND ((asset_type = ETH)) ORDER BY number desc,"accounts"."id" ASC LIMIT 1`).
		WithReply(response)

	isInsertedCalled := false

	assetType := models.Ethereum
	stellarAccountID := "GAJN3BQGETOKZYI6BYNQJZ5ZQWURZX5TXNOEF7UPQAB4BHBCW5JEIOAD"
	defer func() {
		_, err := s.accountService.NewAccount(assetType, stellarAccountID)
		s.NoError(err)
		s.True(isInsertedCalled)
	}()

	mocket.Catcher.NewMock().
		WithQuery(`INSERT  INTO "accounts" ("created_at","updated_at","deleted_at","asset_type","stellar_account_id","number") VALUES (?,?,?,?,?,?)`).
		WithCallback(func(query string, values []driver.NamedValue) {
			isInsertedCalled = true
			var mappedValues []interface{}
			someGeneratedID := 12
			mappedValues = append(mappedValues, someGeneratedID)
			for _, value := range values {
				mappedValues = append(mappedValues, value.Value)
			}

			byteData, err := json.Marshal(mappedValues)
			s.NoError(err)

			createdAccount := models.Account{}
			err = createdAccount.UnmarshalJSON(byteData)
			s.NoError(err)

			s.EqualValues(stellarAccountID, createdAccount.StellarAccountID)
			s.EqualValues(assetType, createdAccount.AssetType)
			s.EqualValues(someGeneratedID, createdAccount.ID)
			s.EqualValues(lastAccountNumber+1, createdAccount.Number)
		})
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
