package db

import (
	"encoding/json"
	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAccount_EthereumAddress(t *testing.T) {
	var accountNumber uint64 = 10
	account := NewAccount(
		Bitcoin,
		"GAKPLJ62YQNPFNMOD4LO5MGD6626VVGWTAEZBLQJYFVV5TCW6OR5ER6D",
		accountNumber)
	mnemonic := "salmon uphold prefer wagon flower attract menu coil confirm humor custom alarm describe deposit family"
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	assert.NoError(t, err)

	assert.Equal(t, account.EthereumAddress(wallet), "0x85C834d966966D8637fe6cd319F9f3B370C1CB03")
}

func TestAccount_UnmarshalJSON(t *testing.T) {
	var accountID uint = 10
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()
	var deletedAt *time.Time
	assetType := Bitcoin
	stellarAccountID := "GBC5IFZM74RSPBN4X5OH35UDWTHROX23YPSK6K5ZS73OEDLGERVZ64AC"
	var accountNumber uint = 20

	fields := []interface{}{
		accountID,
		createdAt,
		updatedAt,
		deletedAt,
		assetType,
		stellarAccountID,
		accountNumber,
	}
	byteData, err := json.Marshal(fields)
	assert.NoError(t, err)

	createdAccount := Account{}
	err = createdAccount.UnmarshalJSON(byteData)
	assert.NoError(t, err)

	assert.EqualValues(t, accountID, createdAccount.ID)
	assert.EqualValues(t, createdAt, createdAccount.CreatedAt)
	assert.EqualValues(t, updatedAt, createdAccount.UpdatedAt)
	assert.EqualValues(t, deletedAt, createdAccount.DeletedAt)
	assert.EqualValues(t, assetType, createdAccount.AssetType)
	assert.EqualValues(t, stellarAccountID, createdAccount.StellarAccountID)
	assert.EqualValues(t, accountNumber, createdAccount.Number)
}
