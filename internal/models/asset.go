package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type AssetType string

const (
	Ethereum AssetType = "ETH"
)

type Asset struct {
	gorm.Model
	AssetType AssetType
	Decimals  uint8
}
