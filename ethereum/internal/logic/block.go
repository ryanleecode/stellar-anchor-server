package logic

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type Block struct {
	number uint64
}

func NewBlock(num uint64) *Block {
	return &Block{
		number: num,
	}
}

func (b Block) Number() uint64 {
	return b.number
}

func (b Block) IsBehind(o Block) bool {
	return b.number < o.number
}

type BlockService struct {
	blockchain     EthereumBlockchain
	canAddToLedger func(tx EthereumTransaction) (bool, error)
}

type EthereumBlockchain interface {
	BlockByNumber(ctx context.Context, num uint64) (Block, error)
	HeadBlockNumber(ctx context.Context) (uint64, error)
	TransactionsFor(ctx context.Context, b Block) ([]EthereumTransaction, error)
}

type AnchorLedger interface {
	LastProcessedBlock() (Block, error)
	AddTx(t EthereumTransaction) error
	AddBlock(b Block) error
}

func NewBlockService(
	b EthereumBlockchain,
	canAddToLedger func(tx EthereumTransaction) (bool, error),
) *BlockService {
	return &BlockService{
		blockchain: b, canAddToLedger: canAddToLedger,
	}
}

func (s BlockService) ProcessNextBlock(ctx context.Context, ledger AnchorLedger) (bool, error) {
	b, err := ledger.LastProcessedBlock()
	if err != nil {
		return false, errors.Wrap(err, "retrieve last processed block failed")
	}

	log.Tracef("the last processed block was %d", b.Number())

	headBNum, err := s.blockchain.HeadBlockNumber(ctx)
	if err != nil {
		return false, errors.Wrap(err, "retrieve head block number failed")
	}
	nxtB := NewBlock(b.number + 1)
	headB := NewBlock(headBNum)
	if !nxtB.IsBehind(*headB) {
		log.Tracef("block %d is not behind head block %d. no blocks processed", b.Number(), headB.Number())
		return false, nil
	}

	log.Tracef("now processing block %d", nxtB.Number())

	err = ledger.AddBlock(*nxtB)
	if err != nil {
		return false, errors.Wrapf(err, "add block %d failed", nxtB.number)
	}
	txs, err := s.blockchain.TransactionsFor(ctx, *nxtB)
	if err != nil {
		return false, errors.Wrapf(err,
			"retrieve transaction for block %d failed", nxtB.number)
	}

	log.Tracef("block %d has %d transactions", nxtB.Number(), len(txs))

	for _, tx := range txs {
		shouldAdd, err := s.canAddToLedger(tx)
		if err != nil {
			return false, errors.Wrapf(err,
				"failed to determine whether tx %s from block %d can be added to the anchor ledger",
				tx.hash, nxtB.Number())
		}
		if !shouldAdd {
			log.Tracef("tx %s cannot be added to our ledger", tx.Hash().Hex())
			continue
		}

		err = ledger.AddTx(tx)
		if err != nil {
			return false, errors.Wrapf(err,
				"failed to add tx %s from block %d to the anchor ledger", tx.hash, nxtB.Number())
		}
	}

	return true, nil
}
