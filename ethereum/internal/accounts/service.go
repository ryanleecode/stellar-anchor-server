package accounts

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/jinzhu/gorm"
	"github.com/stellar/go/support/log"
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

func (s *EthereumAccountService) FindAccount(address string) *EthereumAccount {
	var account AnchorAccount

	if result := s.db.Where(AnchorAccount{Address: address}).First(&account); result.Error != nil {
		log.Warn(result.Error)

		return nil
	}

	return NewEthereumAccount(
		account.StellarAccountID, uint(account.Number), s.wallet)
}

func (s *EthereumAccountService) GetDepositingAccount(stellarAccountID string) (Account, error) {
	var existingAccount AnchorAccount
	if !s.db.
		Table(AnchorAccount{}.TableName()).
		Where(&AnchorAccount{StellarAccountID: stellarAccountID}).
		Scan(existingAccount).RecordNotFound() {
		return NewEthereumAccount(
			stellarAccountID, uint(existingAccount.Number), s.wallet), nil
	}

	var lastAccount AnchorAccount
	s.db.
		Table(AnchorAccount{}.TableName()).
		Order("number desc").
		First(&lastAccount)
	newAccount := NewEthereumAccount(
		stellarAccountID, uint(lastAccount.Number+1), s.wallet)
	ethAccount := AnchorAccount{}.FromEthereumAccount(newAccount)
	s.db.
		Table(AnchorAccount{}.TableName()).
		Create(&ethAccount)

	return newAccount, nil
}
