package mock

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"math/big"
)

type WalletMock struct {
	mock.Mock
}

func (w WalletMock) URL() accounts.URL {
	args := w.Called()

	return args.Get(0).(accounts.URL)
}

func (w WalletMock) Status() (string, error) {
	args := w.Called()

	return args.String(0), args.Error(1)
}

func (w WalletMock) Open(passphrase string) error {
	args := w.Called(passphrase)

	return args.Error(0)
}

func (w WalletMock) Close() error {
	args := w.Called()

	return args.Error(0)
}

func (w WalletMock) Accounts() []accounts.Account {
	args := w.Called()

	return args.Get(0).([]accounts.Account)
}

func (w WalletMock) Contains(account accounts.Account) bool {
	args := w.Called(account)

	return args.Bool(0)
}

func (w WalletMock) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	args := w.Called(path, pin)

	return args.Get(0).(accounts.Account), args.Error(1)
}

func (w WalletMock) SelfDerive(bases []accounts.DerivationPath, chain ethereum.ChainStateReader) {
	w.Called(bases, chain)
}

func (w WalletMock) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	args := w.Called(account, mimeType, data)

	return args.Get(0).([]byte), args.Error(1)
}

func (w WalletMock) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	args := w.Called(account, passphrase, mimeType, data)

	return args.Get(0).([]byte), args.Error(1)
}

func (w WalletMock) SignText(account accounts.Account, text []byte) ([]byte, error) {
	args := w.Called(account, text)

	return args.Get(0).([]byte), args.Error(1)
}

func (w WalletMock) SignTextWithPassphrase(
	account accounts.Account, passphrase string, hash []byte,
) ([]byte, error) {
	args := w.Called(account, passphrase)

	return args.Get(0).([]byte), args.Error(1)
}

func (w WalletMock) SignTx(
	account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {

	args := w.Called(account, tx, chainID)

	return args.Get(0).(*types.Transaction), args.Error(1)

}

func (w WalletMock) SignTxWithPassphrase(
	account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	args := w.Called(account, passphrase, tx, chainID)

	return args.Get(0).(*types.Transaction), args.Error(1)

}
