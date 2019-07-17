package accounts

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/jinzhu/gorm"
	"stellar-fi-anchor/internal/db"
)

type Service struct {
	db     *gorm.DB
	wallet accounts.Wallet
}

func NewService(wallet accounts.Wallet, db *gorm.DB) *Service {
	return &Service{
		wallet: wallet,
		db:     db,
	}
}

func (s *Service) NewAccount(assetType db.AssetType, stellarAccountID string) (*db.Account, error) {
	var lastAccount db.Account
	s.db.Order("number desc").Where("asset_type = ?", assetType).First(&lastAccount)
	newAccount := db.NewAccount(assetType, stellarAccountID, lastAccount.Number+1)
	s.db.Create(&newAccount)

	return newAccount, nil
}
