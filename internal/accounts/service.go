package accounts

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/jinzhu/gorm"
	"stellar-fi-anchor/internal/asset"
)

type Account = interface {
	DepositInstructions() string
}

type EthereumAccountService struct {
	db     *gorm.DB
	wallet accounts.Wallet
}

func NewService(wallet accounts.Wallet, db *gorm.DB) *EthereumAccountService {
	return &EthereumAccountService{
		wallet: wallet,
		db:     db,
	}
}

func (s *EthereumAccountService) CanDeposit(assetType string) bool {
	return assetType == string(asset.Ethereum)
}

func (s *EthereumAccountService) GetDepositingAccount(stellarAccountID string) (Account, error) {
	var existingAccount GenericAccount
	if !s.db.
		Table(GenericAccount{}.TableName()).
		Where(&GenericAccount{AssetType: asset.Ethereum, StellarAccountID: stellarAccountID}).
		Scan(existingAccount).RecordNotFound() {
		return NewEthereumAccount(
			existingAccount.StellarAccountID, existingAccount.Number, s.wallet), nil
	}

	var lastAccount GenericAccount
	s.db.
		Table(GenericAccount{}.TableName()).
		Order("number desc").
		Where(&GenericAccount{AssetType: asset.Ethereum}).
		First(&lastAccount)
	newAccount := NewEthereumAccount(stellarAccountID, lastAccount.Number+1, s.wallet)
	s.db.
		Table(GenericAccount{}.TableName()).
		Create(&newAccount)

	return newAccount, nil
}
