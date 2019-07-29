package data

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/params"

	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"

	"github.com/stellar/go/clients/horizonclient"

	"github.com/stellar/go/keypair"

	"github.com/stellar/go/txnbuild"

	_ "github.com/streadway/amqp"
)

type Issuer struct {
	issuingKP         *keypair.Full
	horizonClient     horizonclient.ClientInterface
	networkPassphrase string
	asset             txnbuild.CreditAsset
}

func NewIssuer(
	issuingKP *keypair.Full,
	horizonClient horizonclient.ClientInterface,
	networkPassphrase string,
	asset txnbuild.CreditAsset,
) *Issuer {
	return &Issuer{
		issuingKP:         issuingKP,
		horizonClient:     horizonClient,
		networkPassphrase: networkPassphrase,
		asset:             asset,
	}
}

func (i *Issuer) IssueWithMemo(destStellarAddr string, amount int64, memo txnbuild.Memo) (issueTxHash string, err error) {
	issuer := i.issuingKP.Address()
	issuerAct, err := i.horizonClient.AccountDetail(
		horizonclient.AccountRequest{
			AccountID: issuer,
		})
	if err != nil {
		return "", errors.Wrapf(err,
			"failed to load issuer account details for %s", issuer)
	}
	destinationAct, err := i.horizonClient.AccountDetail(
		horizonclient.AccountRequest{
			AccountID: destStellarAddr,
		})
	if err != nil {
		return "", errors.Wrapf(err,
			"failed to load destination account details for %s", destStellarAddr)
	}

	availableTrust, err := i.AvailableTrust(i.asset, destinationAct)
	if availableTrust < float64(amount)/params.Ether {
		return "", fmt.Errorf(
			"recipient must have %f trust for the asset issued by %s but instead has %d",
			float64(amount)/params.Ether, issuer, availableTrust)
	}

	// https://www.stellar.org/developers/guides/concepts/assets.html
	formattedAmount := fmt.Sprintf("%.7f", float64(amount)/params.Ether)

	issueAssetTx := txnbuild.Transaction{
		SourceAccount: &issuerAct,
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				Destination: destStellarAddr,
				Amount:      formattedAmount,
				Asset:       i.asset,
			},
		},
		Timebounds: txnbuild.NewInfiniteTimeout(),
		Network:    i.networkPassphrase,
	}
	if memo != nil {
		issueAssetTx.Memo = memo
	}

	issueAssetTxXDR, err := issueAssetTx.BuildSignEncode(i.issuingKP)
	if err != nil {
		return "", errors.Wrapf(
			err, "failed to build, sign, and encode issue asset transaction")
	}
	receipt, err := i.horizonClient.SubmitTransactionXDR(issueAssetTxXDR)
	if err != nil {
		return "", errors.Wrapf(err, "failed to submit issue transaction")
	}

	return receipt.Hash, nil
}

func (i *Issuer) Issue(destinationStellarAddress string, amount int64) (issueTxHash string, err error) {
	return i.IssueWithMemo(destinationStellarAddress, amount, nil)
}

func (i *Issuer) AvailableTrust(asset txnbuild.CreditAsset, account horizon.Account) (float64, error) {
	var availableTrust float64
	for _, b := range account.Balances {
		if b.Asset.Issuer != asset.Issuer || b.Asset.Code != asset.Code {
			continue
		}
		limit, err := strconv.ParseFloat(b.Limit, 64)
		if err != nil {
			return 0, err
		}
		balance, err := strconv.ParseFloat(b.Balance, 64)
		if err != nil {
			return 0, err
		}
		availableTrust = limit - balance
	}

	return availableTrust, nil
}
