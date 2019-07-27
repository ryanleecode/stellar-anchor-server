package accounts

import (
	"encoding/json"
	"fmt"

	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/jinzhu/gorm"
)

type AnchorAccount struct {
	gorm.Model
	StellarAccountID string
	Number           uint64
	Address          string
}

func (AnchorAccount) TableName() string {
	return "accounts"
}

func (a *AnchorAccount) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{
		&a.ID, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt, &a.StellarAccountID, &a.Number, &a.Address}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in AnchorAccount: %d != %d", g, e)
	}
	return nil
}

func (a AnchorAccount) FromEthereumAccount(ethAct *EthereumAccount) *AnchorAccount {
	return &AnchorAccount{
		StellarAccountID: ethAct.StellarAccountID(),
		Number:           uint64(ethAct.Number()),
		Address:          ethAct.Address(),
	}
}

type EthereumWallet interface {
	Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error)
}

type EthereumAccount struct {
	stellarAccountID string
	number           uint
	address          string
}

func NewEthereumAccount(stellarAccountID string, number uint, wallet EthereumWallet) *EthereumAccount {
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", number))
	nativeEthAct, _ := wallet.Derive(path, false)

	return &EthereumAccount{
		stellarAccountID: stellarAccountID,
		number:           number,
		address:          nativeEthAct.Address.Hex(),
	}
}

func (a *EthereumAccount) StellarAccountID() string {
	return a.stellarAccountID
}

func (a *EthereumAccount) Number() uint {
	return a.number
}

func (a *EthereumAccount) Address() string {
	return a.address
}

func (a *EthereumAccount) DepositInstructions() string {
	return a.Address()
}
