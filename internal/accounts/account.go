package accounts

import (
	"encoding/json"
	"fmt"
	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/drdgvhbh/stellar-fi-anchor/internal/asset"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/jinzhu/gorm"
)

type GenericAccount struct {
	gorm.Model
	StellarAccountID string
	AssetType        asset.AssetType
	Number           uint64
}

func (GenericAccount) TableName() string {
	return "accounts"
}

func (a *GenericAccount) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{
		&a.ID, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt, &a.StellarAccountID, &a.AssetType, &a.Number}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Account: %d != %d", g, e)
	}
	return nil
}

type EthereumWallet interface {
	Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error)
}

type EthereumAccount struct {
	StellarAccountID string
	Number           uint64
	wallet           EthereumWallet
}

func NewEthereumAccount(stellarAccountID string, number uint64, wallet EthereumWallet) *EthereumAccount {
	return &EthereumAccount{
		StellarAccountID: stellarAccountID,
		Number:           number,
		wallet:           wallet,
	}
}

func (a *EthereumAccount) Address() string {
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", a.Number))
	nativeEthAct, _ := a.wallet.Derive(path, false)

	return nativeEthAct.Address.Hex()
}

func (a *EthereumAccount) DepositInstructions() string {
	return a.Address()
}
