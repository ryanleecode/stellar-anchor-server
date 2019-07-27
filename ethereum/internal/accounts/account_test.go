package accounts

import (
	"fmt"
	"testing"

	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/mock"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewEthereumAccount(t *testing.T) {
	stellarAccountID := "GDCANY73N6JU5TGXQJI6NSDYKCPPNO5EXZNO4WOCKTQLM7MCMEQKQVFA"

	wallet := mock.WalletMock{}
	actNumber := uint(11)
	ethAddress := "0x87fa4adb38EF0bF29F3EE8035db10fd13B70c0d1"
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", actNumber))
	wallet.On("Derive", path, false).Return(
		accounts.Account{Address: common.HexToAddress(ethAddress)}, nil)

	newAct := NewEthereumAccount(
		stellarAccountID, actNumber, wallet)

	assert.EqualValues(t, ethAddress, newAct.Address())
	assert.EqualValues(t, ethAddress, newAct.DepositInstructions())
	assert.EqualValues(t, stellarAccountID, newAct.StellarAccountID)
	assert.EqualValues(t, actNumber, newAct.Number)
}
