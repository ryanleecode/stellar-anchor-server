package data

import (
	"time"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/logic"

	"github.com/jinzhu/gorm"
)

type LedgerConfig struct {
	BlockTableName       string
	TransactionTableName string
}

type Ledger struct {
	conf LedgerConfig
}

func NewLedger() *Ledger {
	return NewLedgerWithConfig(LedgerConfig{
		BlockTableName:       "blocks",
		TransactionTableName: "transactions",
	})
}

func NewLedgerWithConfig(conf LedgerConfig) *Ledger {
	return &Ledger{
		conf: conf,
	}
}

func (l Ledger) AddTx(transaction logic.EthereumTransaction, db *gorm.DB) error {
	tx := Transaction{
		ETHTxHash:   transaction.Hash().Hex(),
		BlockNumber: transaction.BlockNumber(),
		IsProcessed: false,
	}
	if err := db.Table(l.conf.TransactionTableName).
		Create(&tx).Error; err != nil {
		return err
	}

	return nil
}

func (l Ledger) AddBlock(block logic.Block, db *gorm.DB) error {
	b := NewBlock(block.Number())
	err := db.Table(l.conf.BlockTableName).
		Create(&b).Error
	if err != nil {
		return err
	}

	return nil
}

func (l Ledger) LastProcessedBlock(db *gorm.DB) (logic.Block, error) {
	b := Block{}
	if err := db.Table(l.conf.BlockTableName).
		Order("number desc").
		FirstOrInit(&b, Block{Number: 0}).Error; err != nil {
		return logic.Block{}, err
	}

	return *logic.NewBlock(b.Number), nil
}

func (l Ledger) NextUnprocessed(db *gorm.DB) (*Transaction, error) {
	tx := Transaction{}
	txTable := db.Table(l.conf.TransactionTableName)
	result := txTable.
		Where("is_processed = ?", false).
		First(&tx)
	if result.RecordNotFound() {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &tx, nil
}

func (l Ledger) Process(tx *logic.EthereumTransaction, db *gorm.DB) error {
	dbTx := Transaction{}
	txTable := db.Table(l.conf.TransactionTableName)
	if err := txTable.
		Where("eth_tx_hash = ?", tx.Hash().Hex()).
		First(&dbTx).Error; err != nil {
		return err
	}

	dbTx.IsProcessed = true
	txTable.Save(&dbTx)

	return nil
}

func (l Ledger) Complete(tx *logic.EthereumTransaction, stellarTXHash string, db *gorm.DB) error {
	dbTx := Transaction{}
	txTable := db.Table(l.conf.TransactionTableName)
	if err := txTable.
		Where("eth_tx_hash = ?", tx.Hash().Hex()).
		First(&dbTx).Error; err != nil {
		return err
	}

	dbTx.StellarTxHash = &stellarTXHash
	txTable.Save(&dbTx)

	return nil
}

type Block struct {
	Number    uint64 `gorm:"primary_key;type:bigint"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	tableName string
}

func NewBlock(number uint64) *Block {
	return &Block{
		Number: number,
	}
}

func (b Block) TableName() string {
	return b.tableName
}

type Transaction struct {
	ETHTxHash     string  `gorm:"primary_key"`
	StellarTxHash *string `gorm:"unique_index"`
	IsProcessed   bool    `gorm:"not null"`
	Block         Block   `gorm:"foreignkey:BlockNumber"`
	BlockNumber   uint64  `gorm:"type:bigint"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
	tableName     string
}

func (t Transaction) TableName() string {
	return t.tableName
}
