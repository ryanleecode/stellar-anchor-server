package data

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/logic"
	"github.com/pkg/errors"
)

type LogicGateway struct {
	ledger         *Ledger
	blockchain     *EthereumBlockchain
	db             *gorm.DB
	accountStorage *AccountStorage
}

func NewLogicGateway(
	l *Ledger,
	blockchain *EthereumBlockchain,
	accountStorage *AccountStorage,
	db *gorm.DB,
) *LogicGateway {
	db.AutoMigrate(
		&Block{tableName: l.conf.BlockTableName})
	db.AutoMigrate(
		&Transaction{tableName: l.conf.TransactionTableName}).
		AddForeignKey("block_number", fmt.Sprintf("%s(number)", l.conf.BlockTableName),
			"RESTRICT", "RESTRICT")
	db.AutoMigrate(&Account{tableName: accountStorage.conf.AccountsTableName})
	return &LogicGateway{
		ledger:         l,
		blockchain:     blockchain,
		accountStorage: accountStorage,
		db:             db,
	}
}

func (g LogicGateway) BlockByNumber(ctx context.Context, num uint64) (logic.Block, error) {
	return g.blockchain.BlockByNumber(ctx, num)
}

func (g LogicGateway) HeadBlockNumber(ctx context.Context) (uint64, error) {
	return g.blockchain.HeadBlockNumber(ctx)
}

func (g LogicGateway) TransactionsFor(ctx context.Context, b logic.Block) ([]logic.EthereumTransaction, error) {
	return g.blockchain.TransactionsFor(ctx, b)
}

func (g LogicGateway) Begin() *LogicGateway {
	return &LogicGateway{ledger: g.ledger, blockchain: g.blockchain, db: g.db.Begin()}
}

func (g LogicGateway) Commit() *LogicGateway {
	return &LogicGateway{ledger: g.ledger, blockchain: g.blockchain, db: g.db.Commit()}
}

func (g LogicGateway) RollbackUnlessCommitted() *LogicGateway {
	return &LogicGateway{ledger: g.ledger, blockchain: g.blockchain, db: g.db.RollbackUnlessCommitted()}
}

func (g LogicGateway) Rollback() *LogicGateway {
	return &LogicGateway{ledger: g.ledger, blockchain: g.blockchain, db: g.db.Rollback()}
}

func (g LogicGateway) AddTx(transaction logic.EthereumTransaction) error {
	return g.ledger.AddTx(transaction, g.db)
}

func (g LogicGateway) AddBlock(block logic.Block) error {
	return g.ledger.AddBlock(block, g.db)
}

func (g LogicGateway) LastProcessedBlock() (logic.Block, error) {
	return g.ledger.LastProcessedBlock(g.db)
}

func (g LogicGateway) NextUnprocessed(ctx context.Context) (*logic.EthereumTransaction, error) {
	dbTx, err := g.ledger.NextUnprocessed(g.db)
	if err != nil {
		return nil, errors.Wrap(err, "cannot retrieve next unprocessed tx from db")
	}
	if dbTx == nil {
		return nil, nil
	}

	tx, err := g.blockchain.TransactionByHash(ctx, dbTx.ETHTxHash)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find tx by hash %s", dbTx.ETHTxHash)
	}

	return tx, nil
}

func (g LogicGateway) Process(tx *logic.EthereumTransaction) error {
	return g.ledger.Process(tx, g.db)
}

func (g LogicGateway) Complete(tx *logic.EthereumTransaction, stellarTXHash string) error {
	return g.ledger.Complete(tx, stellarTXHash, g.db)
}

func (g LogicGateway) FindByStellar(addr string) (*logic.DepositAccount, error) {
	return g.accountStorage.FindByStellar(addr, g.db)
}
func (g LogicGateway) FindByEth(addr string) (*logic.DepositAccount, error) {
	return g.accountStorage.FindByEth(addr, g.db)
}

func (g LogicGateway) New(stellarAcctAddr string) (*logic.DepositAccount, error) {
	return g.accountStorage.New(stellarAcctAddr, g.db)
}
