package models

import (
	"encoding/json"
	"fmt"
	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	AssetType        AssetType
	StellarAccountID string
	Number           uint64
}

func NewAccount(assetType AssetType, stellarAccountID string, number uint64) *Account {
	return &Account{
		AssetType:        assetType,
		StellarAccountID: stellarAccountID,
		Number:           number,
	}
}

func (a *Account) EthereumAddress(wallet accounts.Wallet) string {
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", a.Number))

	ethAccount, _ := wallet.Derive(path, false)

	return ethAccount.Address.Hex()
}
func (a *Account) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{
		&a.ID, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt, &a.AssetType, &a.StellarAccountID, &a.Number}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Account: %d != %d", g, e)
	}
	return nil
}
