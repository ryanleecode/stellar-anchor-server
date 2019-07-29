package mock

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
)

type EthClient struct {
	mock.Mock
}

func (c *EthClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	args := c.Called(ctx, number)

	return args.Get(0).(*types.Block), args.Error(1)
}

func (c *EthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	args := c.Called(ctx, number)

	return args.Get(0).(*types.Header), args.Error(1)
}

func (c *EthClient) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	args := c.Called(ctx, hash)

	return args.Get(0).(*types.Transaction), args.Bool(1), args.Error(2)
}

func (c *EthClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	args := c.Called(ctx, txHash)

	return args.Get(0).(*types.Receipt), args.Error(2)
}
