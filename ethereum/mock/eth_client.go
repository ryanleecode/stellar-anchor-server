package mock

import (
	"context"
	"math/big"

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
