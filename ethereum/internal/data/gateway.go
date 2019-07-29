package data

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

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

type LogicGatewayTracer struct {
	logger  log.Ext1FieldLogger
	gateway Gateway
}

type Gateway interface {
	logic.DepositAccountStorage
	logic.EthereumBlockchain
	logic.Deposits
	logic.AnchorLedger
}

func NewLogicGatewayTracer(
	logger log.Ext1FieldLogger,
	gateway Gateway,
) *LogicGatewayTracer {
	return &LogicGatewayTracer{
		logger:  logger,
		gateway: gateway,
	}
}

func (t LogicGatewayTracer) BlockByNumber(ctx context.Context, num uint64) (logic.Block, error) {
	return t.gateway.BlockByNumber(ctx, num)
}

func (t LogicGatewayTracer) HeadBlockNumber(ctx context.Context) (uint64, error) {
	num, err := t.gateway.HeadBlockNumber(ctx)
	t.logger.WithField("head_block_number", num).
		WithField("err", err).
		Traceln("head block number")

	return num, err
}

func (t LogicGatewayTracer) TransactionsFor(ctx context.Context, b logic.Block) ([]logic.EthereumTransaction, error) {
	t.logger.WithField("block_number", b.Number()).Traceln("transactions for block")
	txs, err := t.gateway.TransactionsFor(ctx, b)
	if err == nil {
		t.logger.WithFields(log.Fields{
			"block_number": b.Number(),
			"tx_count":     len(txs),
		}).Traceln("transactions for block found")
	}

	return txs, err
}

func (t LogicGatewayTracer) AddTx(tx logic.EthereumTransaction) error {
	t.logger.WithFields(tx.Fields()).Traceln("add transaction")
	return t.gateway.AddTx(tx)
}

func (t LogicGatewayTracer) AddBlock(block logic.Block) error {
	t.logger.WithFields(log.Fields{
		"number": block.Number(),
	}).Traceln("add block")
	return t.gateway.AddBlock(block)
}

func (t LogicGatewayTracer) LastProcessedBlock() (logic.Block, error) {
	b, err := t.gateway.LastProcessedBlock()
	if err == nil {
		t.logger.WithFields(log.Fields{
			"number": b.Number(),
		}).Traceln("last processed block")
	}

	return b, err
}

func (t LogicGatewayTracer) NextUnprocessed(ctx context.Context) (*logic.EthereumTransaction, error) {
	tx, err := t.gateway.NextUnprocessed(ctx)
	if err == nil && tx != nil {
		t.logger.WithFields(tx.Fields()).Debugln("next unprocessed transaction")
	}

	return tx, err
}

func (t LogicGatewayTracer) Process(tx *logic.EthereumTransaction) error {
	t.logger.WithFields(tx.Fields()).Traceln("process transaction")
	return t.gateway.Process(tx)
}

func (t LogicGatewayTracer) Complete(tx *logic.EthereumTransaction, stellarTXHash string) error {
	t.logger.WithFields(tx.Fields()).
		WithField("stellar_tx_hash", stellarTXHash).
		Traceln("complete transaction")
	return t.gateway.Complete(tx, stellarTXHash)
}

func (t LogicGatewayTracer) FindByStellar(addr string) (*logic.DepositAccount, error) {
	t.logger.
		WithField("addr", addr).
		Traceln("find by stellar")
	return t.gateway.FindByStellar(addr)
}

func (t LogicGatewayTracer) FindByEth(addr string) (*logic.DepositAccount, error) {
	t.logger.
		WithField("addr", addr).
		Traceln("find by eth")
	return t.gateway.FindByEth(addr)
}

func (t LogicGatewayTracer) New(stellarAcctAddr string) (*logic.DepositAccount, error) {
	t.logger.
		WithField("stellar_acct_addr", stellarAcctAddr).
		Traceln("new")
	return t.gateway.New(stellarAcctAddr)
}
