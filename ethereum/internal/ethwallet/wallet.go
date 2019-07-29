package ethwallet

import (
	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/data"

	"github.com/pkg/errors"

	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/ethereum/go-ethereum/accounts"
)

type Wallet struct {
	wallet *hdwallet.Wallet
}

func NewWallet(wallet *hdwallet.Wallet) *Wallet {
	return &Wallet{
		wallet: wallet,
	}
}

func (w Wallet) PrivateKeyBytes(path accounts.DerivationPath) ([]byte, error) {
	acct, err := w.wallet.Derive(path, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not derive eth account")
	}
	pk, err := w.wallet.PrivateKeyBytes(acct)
	if err != nil {
		return nil, errors.Wrapf(err, "could not derive private key for %s", acct.Address)
	}
	return pk, nil
}

func (w Wallet) Derive(path accounts.DerivationPath, pin bool) (*data.EthereumAccount, error) {
	acct, err := w.wallet.Derive(path, false)
	if err != nil {
		return nil, errors.Wrap(err, "could not derive eth account")
	}

	return data.NewEthereumAccount(acct.Address.Hex()), nil
}
