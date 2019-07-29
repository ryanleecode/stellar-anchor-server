package data_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/data"

	mocks "github.com/drdgvhbh/stellar-fi-anchor/ethereum/mock"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/logic"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/suite"
)

type MockBlockCtor struct {
	mock.Mock
}

func (m *MockBlockCtor) NewBlock(num uint64) logic.Block {
	args := m.Called(num)

	return args.Get(0).(logic.Block)
}

type EthereumBlockchainSuite struct {
	suite.Suite
	blockCtor MockBlockCtor
	ethClient mocks.EthClient
	ctx       context.Context
}

func (s *EthereumBlockchainSuite) SetupSuite() {
}

func (s *EthereumBlockchainSuite) SetupTest() {
	s.blockCtor = MockBlockCtor{}
	s.ethClient = mocks.EthClient{}
	s.ctx = context.TODO()
}

func (s *EthereumBlockchainSuite) TestBlockByNumber_FindsTheBlock() {
	num := uint64(10)
	s.ethClient.On("BlockByNumber", s.ctx, new(big.Int).SetUint64(num)).Return(
		types.NewBlock(&types.Header{Number: big.NewInt(int64(num))},
			nil, nil, nil),
		nil)

	expectedB := logic.NewBlock(num)
	s.blockCtor.On("NewBlock", num).Return(*expectedB)

	bc := data.NewEthereumBlockchain(&s.ethClient, s.blockCtor.NewBlock)
	b, err := bc.BlockByNumber(s.ctx, num)

	s.EqualValues(*expectedB, b)
	s.NoError(err)
	s.blockCtor.AssertExpectations(s.T())
	s.ethClient.AssertExpectations(s.T())
}

func (s *EthereumBlockchainSuite) TestBlockByNumber_FailsIfEthClientFails() {
	num := uint64(10)

	s.ethClient.On("BlockByNumber", s.ctx, new(big.Int).SetUint64(num)).Return(
		(*types.Block)(nil),
		errors.New(""))

	bc := data.NewEthereumBlockchain(&s.ethClient, s.blockCtor.NewBlock)
	_, err := bc.BlockByNumber(s.ctx, num)

	s.Error(err)
}

func (s *EthereumBlockchainSuite) TestHeadBlockNumber_FindsTheNumber() {
	num := uint64(10)

	s.ethClient.On("HeaderByNumber", s.ctx, (*big.Int)(nil)).Return(
		&types.Header{Number: big.NewInt(int64(num))},
		nil)

	bc := data.NewEthereumBlockchain(&s.ethClient, s.blockCtor.NewBlock)
	n, err := bc.HeadBlockNumber(s.ctx)

	s.EqualValues(num, n)
	s.NoError(err)
}

func (s *EthereumBlockchainSuite) TestHeadBlockNumber_FailsIfEthClientFails() {
	s.ethClient.On("HeaderByNumber", s.ctx, (*big.Int)(nil)).Return(
		(*types.Header)(nil),
		errors.New(""))

	bc := data.NewEthereumBlockchain(&s.ethClient, s.blockCtor.NewBlock)
	_, err := bc.HeadBlockNumber(s.ctx)

	s.Error(err)
}

func TestEthereumBlockchainSuite(t *testing.T) {
	suite.Run(t, new(EthereumBlockchainSuite))
}
