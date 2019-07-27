package internal

import (
	"time"

	"github.com/jinzhu/gorm"
)

type DBConfig struct {
	BlockTableName       string
	TransactionTableName string
}

type DB struct {
	db   *gorm.DB
	conf DBConfig
}

func NewDB(db *gorm.DB) *DB {
	return NewDBWithConfig(db, DBConfig{
		BlockTableName:       "blocks",
		TransactionTableName: "transactions",
	})
}

func NewDBWithConfig(db *gorm.DB, conf DBConfig) *DB {
	db.AutoMigrate(
		Block{tableName: conf.BlockTableName},
		Transaction{tableName: conf.TransactionTableName})
	return &DB{
		db:   db,
		conf: conf,
	}
}

func (d DB) Begin() *DB {
	return NewDBWithConfig(d.db.Begin(), d.conf)
}

func (d DB) RollbackUnlessCommitted() *DB {
	return NewDBWithConfig(d.db.RollbackUnlessCommitted(), d.conf)
}

func (d DB) Commit() *DB {
	return NewDBWithConfig(d.db.Commit(), d.conf)
}

func (d DB) AddBlock(b Block) error {
	if err := d.db.Table(d.conf.BlockTableName).
		Debug().
		Create(&b).Error; err != nil {
		return err
	}

	return nil
}

func (d DB) LastProcessedBlock() (Block, error) {
	b := Block{}
	if err := d.db.Table(d.conf.BlockTableName).
		Debug().
		Order("number desc").
		FirstOrInit(&b, Block{Number: 0}).Error; err != nil {
		return Block{}, err
	}

	return b, nil
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
	ETHTxHash     string `gorm:"primary_key"`
	StellarTxHash string `gorm:"unique_index"`
	IsProcessed   bool   `gorm:"not null"`
	tableName     string
}

func (t Transaction) TableName() string {
	return t.tableName
}
