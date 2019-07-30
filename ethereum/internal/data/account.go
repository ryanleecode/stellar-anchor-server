package data

import (
	"fmt"
	"time"

	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/ethereum/go-ethereum/accounts"

	"github.com/pkg/errors"

	"github.com/drdgvhbh/stellar-anchor-server/ethereum/internal/logic"

	"github.com/jinzhu/gorm"
)

type Account struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
	StellarAccountID string     `gorm:"primary_key;"`
	Number           uint64     `gorm:"unique;not null"`
	Address          string     `gorm:"unique;not null"`
	tableName        string
}

func (a Account) TableName() string {
	return a.tableName
}

type AccountStorageConf struct {
	AccountsTableName string
}

type EthereumAccount struct {
	address string
}

func NewEthereumAccount(address string) *EthereumAccount {
	return &EthereumAccount{address: address}
}

type EthereumWallet interface {
	Derive(path accounts.DerivationPath, pin bool) (*EthereumAccount, error)
}

type AccountStorage struct {
	conf   AccountStorageConf
	wallet EthereumWallet
}

func NewAccountStorage(w EthereumWallet) *AccountStorage {
	return NewAccountStorageWithConf(w, AccountStorageConf{AccountsTableName: "accounts"})
}

func NewAccountStorageWithConf(w EthereumWallet, conf AccountStorageConf) *AccountStorage {
	return &AccountStorage{conf: conf, wallet: w}
}

func (s AccountStorage) FindByStellar(addr string, db *gorm.DB) (*logic.DepositAccount, error) {
	acct := Account{}
	table := db.Table(s.conf.AccountsTableName)
	result := table.Where("stellar_account_id = ?", addr).
		First(&acct)
	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error,
			"could not lookup account with address %s from db", addr)
	}

	return logic.NewDepositAccount(acct.Address, acct.StellarAccountID), nil
}

func (s AccountStorage) FindByEth(addr string, db *gorm.DB) (*logic.DepositAccount, error) {
	acct := Account{}
	table := db.Table(s.conf.AccountsTableName)
	result := table.Where("address = ?", addr).First(&acct)
	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, errors.Wrapf(result.Error,
			"could not lookup account with address %s from db", addr)
	}

	return logic.NewDepositAccount(acct.Address, acct.StellarAccountID), nil
}

func (s AccountStorage) New(stellarAcctAddr string, db *gorm.DB) (*logic.DepositAccount, error) {
	mostRecentAccount := Account{}
	table := db.Table(s.conf.AccountsTableName)
	table.Order("number desc").First(&mostRecentAccount)

	acctNumber := mostRecentAccount.Number + 1
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", acctNumber))
	nativeEthAcct, err := s.wallet.Derive(path, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to derive eth account number %d", acctNumber)
	}

	newAcct := Account{
		StellarAccountID: stellarAcctAddr,
		Number:           acctNumber,
		Address:          nativeEthAcct.address,
	}

	if err := table.Create(&newAcct).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to create new account for %s in database", stellarAcctAddr)
	}

	return logic.NewDepositAccount(newAcct.Address, stellarAcctAddr), nil
}
