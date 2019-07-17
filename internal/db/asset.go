package db

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type AssetType string

const (
	Ethereum AssetType = "ETH"
	Bitcoin  AssetType = "BTC"
)

type Asset struct {
	AssetType AssetType `gorm:"primary_key"`
	Decimals  uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
