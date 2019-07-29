package logic

import (
	"context"
	"math/big"

	log "github.com/sirupsen/logrus"

	"github.com/stellar/go/txnbuild"

	"github.com/pkg/errors"
)

type Hash interface {
	Hex() string
	Bytes() []byte
}

type EthereumTransaction struct {
	from        string
	to          string
	gwei        big.Int
	hash        Hash
	blockNumber uint64
}

func (t EthereumTransaction) From() (ethereumAddress string) {
	return t.from
}

func (t EthereumTransaction) To() (ethereumAddress string) {
	return t.to
}

func (t EthereumTransaction) Gwei() big.Int {
	return t.gwei
}

func (t EthereumTransaction) Hash() (ethereumHash Hash) {
	return t.hash
}

func (t EthereumTransaction) BlockNumber() uint64 {
	return t.blockNumber
}

func NewEthereumTransaction(
	From string,
	To string,
	Gwei big.Int,
	Hash Hash,
	BlockNumber uint64,
) *EthereumTransaction {
	return &EthereumTransaction{
		from:        From,
		to:          To,
		gwei:        Gwei,
		hash:        Hash,
		blockNumber: BlockNumber,
	}
}

type Deposits interface {
	NextUnprocessed(ctx context.Context) (*EthereumTransaction, error)
	Process(tx *EthereumTransaction) error
	Complete(tx *EthereumTransaction, stellarTXHash string) error
}

type Issuer interface {
	IssueWithMemo(destStellarAddr string, amount int64, memo txnbuild.Memo) (issueTxHash string, err error)
}

type AccountDictionary interface {
	FindAccountFrom(ethAddr string) (*StellarAccount, error)
}

type PrecisionPolicy func(number big.Int) int64
type TransactionService struct {
	deposits        Deposits
	dict            AccountDictionary
	precisionPolicy PrecisionPolicy
}

func NewTransactionService(dict AccountDictionary, p PrecisionPolicy, d Deposits) *TransactionService {
	return &TransactionService{
		deposits:        d,
		dict:            dict,
		precisionPolicy: p,
	}
}

func (s TransactionService) ProcessDeposit(ctx context.Context, d Deposits) error {
	tx, err := s.deposits.NextUnprocessed(ctx)
	if err != nil {
		return errors.Wrap(err, "retrieve next unprocessed transaction failed")
	}
	if tx == nil {
		return nil
	}

	acct, err := s.dict.FindAccountFrom(tx.To())
	if err != nil {
		return errors.Wrapf(err, "find stellar account for eth address %s failed", tx.To())
	}

	if err = s.deposits.Process(tx); err != nil {
		return errors.Wrapf(err, "cannot mark eth tx %s as processed", tx.Hash().Hex())
	}
	amount := s.precisionPolicy(tx.Gwei())
	hash, err := acct.Issue(amount, tx)
	if err != nil {
		return errors.Wrapf(err,
			"cannot issue asset for deposit. eth tx: %s. recipient stellar acct: %s", tx.Hash().Hex(), acct.address)
	}
	err = s.deposits.Complete(tx, hash)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"eth_tx_hash":               tx.hash.Hex(),
			"stellar_tx_hash":           hash,
			"amount":                    amount,
			"recipient_stellar_address": acct.address,
		}).Warn("the asset was issued successfully but could not be marked completed in our records")
	}

	return nil
}
