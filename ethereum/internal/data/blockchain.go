package data

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/logic"
	"github.com/pkg/errors"
)

type EthClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
}

type NewLogicBlock func(num uint64) logic.Block

type EthereumBlockchain struct {
	ethClient EthClient
	newBlock  NewLogicBlock
}

func NewEthereumBlockchain(client EthClient, newBlock NewLogicBlock) *EthereumBlockchain {
	return &EthereumBlockchain{ethClient: client, newBlock: newBlock}
}

func (chain EthereumBlockchain) BlockByNumber(ctx context.Context, num uint64) (logic.Block, error) {
	b, err := chain.blockByNumber(ctx, num)
	if err != nil {
		return logic.Block{}, errors.Wrapf(err, "retrieve block %d from ethereum blockchain failed", num)
	}

	return chain.newBlock(b.NumberU64()), nil
}

func (chain EthereumBlockchain) HeadBlockNumber(ctx context.Context) (uint64, error) {
	headHeader, err := chain.ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, errors.Wrapf(err, "retrieve head header from ethereum blockchain failed")
	}

	return headHeader.Number.Uint64(), nil
}

func (chain EthereumBlockchain) TransactionsFor(ctx context.Context, b logic.Block) ([]logic.EthereumTransaction, error) {
	ethBlock, err := chain.blockByNumber(ctx, b.Number())
	if err != nil {
		return nil, errors.Wrapf(err, "retrieve transactions for block %d from ethereum blockchain failed", b.Number())
	}

	var txs []logic.EthereumTransaction
	for _, tx := range ethBlock.Transactions() {
		if tx.To() == nil {
			continue
		}
		lTx, err := chain.logicTx(tx, b.Number())
		if err != nil {
			return nil, errors.Wrapf(err, "cannot convert tx %s to logic eth tx", tx.Hash())
		}

		txs = append(
			txs,
			*lTx)
	}

	return txs, nil
}

func (chain EthereumBlockchain) TransactionByHash(ctx context.Context, hash string) (*logic.EthereumTransaction, error) {
	tx, _, err := chain.ethClient.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find tx %s from the eth blockchain", hash)
	}
	txReceipt, err := chain.ethClient.TransactionReceipt(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, errors.Wrapf(err, "tx %s has no tx receipt.", hash)
	}

	lTx, err := chain.logicTx(tx, txReceipt.BlockNumber.Uint64())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot convert tx %s to logic eth tx", tx.Hash())
	}

	return lTx, nil
}

func (chain EthereumBlockchain) blockByNumber(ctx context.Context, num uint64) (*types.Block, error) {
	return chain.ethClient.BlockByNumber(ctx, new(big.Int).SetUint64(num))
}

func (chain EthereumBlockchain) logicTx(tx *types.Transaction, bNum uint64) (*logic.EthereumTransaction, error) {
	signer := types.NewEIP155Signer(tx.ChainId())
	sender, err := signer.Sender(tx)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find sender of ethereum transaction %s", tx.Hash())
	}

	return logic.NewEthereumTransaction(
		sender.Hex(),
		tx.To().Hex(),
		*tx.Value(),
		tx.Hash(),
		bNum), nil
}
