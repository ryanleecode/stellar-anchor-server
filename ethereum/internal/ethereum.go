package internal

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/params"

	"github.com/pkg/errors"
	"github.com/stellar/go/protocols/horizon"

	"github.com/stellar/go/clients/horizonclient"

	"github.com/stellar/go/keypair"

	"github.com/stellar/go/txnbuild"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/accounts"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	_ "github.com/streadway/amqp"
)

type AccountsService interface {
	FindAccount(address string) *accounts.EthereumAccount
}

type DepositCallback func(eth *accounts.EthereumAccount, txHash string, gwei *big.Int)

type EthereumDepositService struct {
	ethClient  *ethclient.Client
	actService AccountsService
}

func watchForEthereumDeposits(
	sub ethereum.Subscription,
	headers chan *types.Header,
	ethClient *ethclient.Client,
	service AccountsService,
	onDeposit DepositCallback,
) {
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := ethClient.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				// TODO: Add to a failure queue somewhere
				log.Fatal(err)
			}

			for _, tx := range block.Transactions() {
				if tx.To() == nil || tx.Value() == nil {
					continue
				}

				accountAddress := tx.To().Hex()
				ethAccount := service.FindAccount(accountAddress)
				if ethAccount != nil {
					txHash := tx.Hash().Hex()
					gwei := tx.Value()

					onDeposit(ethAccount, txHash, gwei)
				}
			}
		}
	}
}

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

func (i *Issuer) IssueWithMemo(destinationStellarAddress string, amount int64, memo txnbuild.Memo) error {
	issuer := i.issuingKP.Address()
	issuerAct, err := i.horizonClient.AccountDetail(
		horizonclient.AccountRequest{
			AccountID: issuer,
		})
	if err != nil {
		return errors.Wrapf(err, "failed to load account details for %s", issuer)
	}
	destinationAct, err := i.horizonClient.AccountDetail(
		horizonclient.AccountRequest{
			AccountID: destinationStellarAddress,
		})
	if err != nil {
		return errors.Wrapf(err, "failed to load account details for %s", destinationStellarAddress)
	}

	availableTrust, err := i.AvailableTrust(i.asset, destinationAct)
	if availableTrust < amount {
		return errors.New("recipient doesnt not have enough available trust for the asset ")
	}

	// https://www.stellar.org/developers/guides/concepts/assets.html
	formattedAmount := fmt.Sprintf("%.7f", float64(amount)/params.Ether)

	log.Println(formattedAmount, float64(amount), params.Ether)

	issueAssetTx := txnbuild.Transaction{
		SourceAccount: &issuerAct,
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				Destination: destinationStellarAddress,
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
		return errors.Wrapf(err, "failed to build, sign, and encode issue asset transaction")
	}
	_, err = i.horizonClient.SubmitTransactionXDR(issueAssetTxXDR)
	if err != nil {
		return errors.Wrapf(err, "failed to submit issue transaction")
	}

	return nil
}

func (i *Issuer) Issue(destinationStellarAddress string, amount int64) error {
	return i.IssueWithMemo(destinationStellarAddress, amount, nil)
}

func (i *Issuer) AvailableTrust(asset txnbuild.CreditAsset, account horizon.Account) (int64, error) {
	var availableTrust int64
	for _, b := range account.Balances {
		if b.Asset.Issuer != asset.Issuer || b.Asset.Code != asset.Code {
			continue
		}
		limit, err := strconv.ParseInt(i.removeDecimals(b.Limit), 10, 64)
		if err != nil {
			return 0, err
		}
		balance, err := strconv.ParseInt(i.removeDecimals(b.Balance), 10, 64)
		if err != nil {
			return 0, err
		}
		availableTrust = limit - balance
	}

	return availableTrust, nil
}

func (i *Issuer) removeDecimals(number string) string {
	return strings.ReplaceAll(number, ".", "")
}
