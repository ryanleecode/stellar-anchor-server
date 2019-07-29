package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/ethwallet"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/stellar"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/logic"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/data"

	log "github.com/sirupsen/logrus"

	"github.com/stellar/go/txnbuild"

	"github.com/stellar/go/clients/horizonclient"

	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"github.com/stellar/go/keypair"
)

type BootstrapParams interface {
	NetworkPassphrase() string
	Mnemonic() string
	DB() *gorm.DB
	RPCClient() *rpc.Client
}

func Bootstrap(params BootstrapParams, logger *log.Logger) http.Handler {
	w, err := hdwallet.NewFromMnemonic(params.Mnemonic())
	if err != nil {
		logger.Fatalln(errors.Wrap(err, "cannot create hd wallet").Error())
	}
	wallet := ethwallet.NewWallet(w)

	ethClient := ethclient.NewClient(params.RPCClient())
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/%d", 1))
	pk, err := wallet.PrivateKeyBytes(path)
	if err != nil {
		logger.Fatalln(err)
	}
	var pk32 [32]byte
	copy(pk32[:], pk)
	issuingKP, err := keypair.FromRawSeed(pk32)
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Println(issuingKP.Address(), issuingKP.Seed())

	issuer := data.NewIssuer(
		issuingKP,
		horizonclient.DefaultTestNetClient,
		params.NetworkPassphrase(),
		txnbuild.CreditAsset{Code: "ETH", Issuer: issuingKP.Address()})

	db := params.DB()

	ethBlockchain := data.NewEthereumBlockchain(ethClient, func(num uint64) logic.Block {
		return *logic.NewBlock(num)
	})
	ledger := data.NewLedger()
	acctStorage := data.NewAccountStorage(wallet)
	gateway := data.NewLogicGateway(ledger, ethBlockchain, acctStorage, db)
	acctService := logic.NewAccountService(data.NewLogicGatewayTracer(logger, gateway), issuer)

	blockService := logic.NewBlockService(data.NewLogicGatewayTracer(logger, gateway),
		func(tx logic.EthereumTransaction) (bool, error) {
			act, err := acctService.FindAccountFrom(tx.To())
			if err != nil {
				return false, err
			}
			return act != nil, err
		})
	transactionService := logic.NewTransactionService(
		acctService, stellar.FormatToAssetPrecision, data.NewLogicGatewayTracer(logger, gateway))

	c := cron.New()
	if err := c.AddFunc("@every 10s", func() {
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*9900)
		for {
			if ctx.Err() != nil {
				break
			}

			l := gateway.Begin()

			didProcess, err := blockService.ProcessNextBlock(ctx, data.NewLogicGatewayTracer(logger, l))
			if err != nil {
				l.Rollback()
				logger.WithError(err).Error("failed to process next block")
				break
			}
			if !didProcess {
				l.Rollback()
				break
			}
			l.Commit()
		}
	}); err != nil {
		logger.Fatalln(errors.Wrapf(err, "failed to add process blocks cron job"))
	}
	if err := c.AddFunc("@every 5s", func() {
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*4900)
		for {
			if ctx.Err() != nil {
				break
			}

			l := gateway.Begin()

			err := transactionService.ProcessDeposit(ctx, data.NewLogicGatewayTracer(logger, l))
			if err != nil {
				l.Rollback()
				logger.WithError(err).Error("failed to process next deposit transaction")

				errCause := errors.Cause(err)
				switch errCause.(type) {
				case *horizonclient.Error:
					horizonErr := errCause.(*horizonclient.Error)
					logger.WithError(errCause).WithFields(log.Fields{
						"title":  horizonErr.Problem.Title,
						"type":   horizonErr.Problem.Type,
						"detail": horizonErr.Problem.Detail,
						"status": horizonErr.Problem.Status,
						"extras": horizonErr.Problem.Extras,
					}).Errorf("horizon error")
					break
				default:
					break
				}
				break
			}
			l.Commit()
		}
	}); err != nil {
		logger.Fatalln(errors.Wrapf(err, "failed to add process transactions cron job"))
	}
	c.Start()

	return NewRootHandler(acctService)
}
