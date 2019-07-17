package accounts

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/jinzhu/gorm"
	"stellar-fi-anchor/internal/models"
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

func (s *Service) NewAccount(assetType models.AssetType, stellarAccountID string) (*models.Account, error) {
	var lastAccount models.Account
	s.db.Order("number desc").Where("asset_type = ?", assetType).First(&lastAccount)
	newAccount := models.NewAccount(assetType, stellarAccountID, lastAccount.Number+1)
	s.db.Create(&newAccount)

	return newAccount, nil
}
