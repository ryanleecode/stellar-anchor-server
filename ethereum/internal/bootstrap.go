package internal

import (
	"context"
	"math"
	"math/big"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/stellar/go/txnbuild"

	"github.com/stellar/go/clients/horizonclient"

	hdwallet "github.com/drdgvhbh/go-ethereum-hdwallet"
	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
)

type BootstrapParams interface {
	NetworkPassphrase() string
	Mnemonic() string
	DB() *gorm.DB
	RPCClient() *rpc.Client
}

func Bootstrap(params BootstrapParams) http.Handler {
	params.DB().AutoMigrate(accounts.AnchorAccount{})

	wallet, err := hdwallet.NewFromMnemonic(params.Mnemonic())
	if err != nil {
		log.Fatalln(errors.Wrap(err, "cannot create hd wallet").Error())
	}

	ethClient := ethclient.NewClient(params.RPCClient())
	headers := make(chan *types.Header)
	sub, err := ethClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalln(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	issuingAct, err := wallet.Derive(path, false)
	pk, err := wallet.PrivateKeyBytes(issuingAct)
	if err != nil {
		log.Fatalln(err)
	}
	var pk32 [32]byte
	copy(pk32[:], pk)
	issuingKP, err := keypair.FromRawSeed(pk32)
	if err != nil {
		log.Fatalln(err)
	}

	actService := accounts.NewService(wallet, params.DB())
	issuer := NewIssuer(
		issuingKP,
		horizonclient.DefaultTestNetClient,
		params.NetworkPassphrase(),
		txnbuild.CreditAsset{Code: "ETH", Issuer: issuingKP.Address()})

	go func() {
		for {
			db := NewDB(params.DB().Begin())

			func() {
				defer db.RollbackUnlessCommitted()

				block, err := db.LastProcessedBlock()
				if err != nil {
					log.WithError(err).Warn("failed to retrieve last processed block")
					return
				}
				headHeader, err := ethClient.HeaderByNumber(context.Background(), nil)
				if err != nil {
					log.WithError(err).Warn("failed to retrieve head block header")
					return
				}

				var chainBlock *types.Block
				var blockNumber uint64
				isOutOfSync := block.Number < (headHeader.Number.Uint64() - 1)
				if isOutOfSync {
					nextNum := block.Number + 1
					chainBlock, err = ethClient.BlockByNumber(context.Background(), new(big.Int).SetUint64(nextNum))
					if err != nil {
						log.WithError(err).
							WithField("block_number", nextNum).
							Warnf("failed to retrieve block")
						return
					}
					blockNumber = nextNum
				} else {
					select {
					case err := <-sub.Err():
						log.Fatal(err)
					case header := <-headers:
						chainBlock, err = ethClient.BlockByHash(context.Background(), header.Hash())
						if err != nil {
							log.WithError(err).Warn("failed to retrieve head block")
							return
						}
						blockNumber = chainBlock.Number().Uint64()
					}
				}

				for _, tx := range chainBlock.Transactions() {
					notAValueTx := tx.To() == nil || tx.Value() == nil
					if notAValueTx {
						continue
					}

					destAddr := tx.To().Hex()
					account := actService.FindAccount(destAddr)
					acctNotInOurRecords := account == nil
					if acctNotInOurRecords {
						continue
					}

					txHash := tx.Hash().Hex()
					gwei := tx.Value()
					depositAmount := int64(truncateToStellarPrecision(gwei))
					err := issuer.IssueWithMemo(
						account.StellarAccountID(), depositAmount, txnbuild.MemoText(txHash))
					if err != nil {
						errCause := errors.Cause(err)
						switch errCause.(type) {
						case *horizonclient.Error:
							log.Println(errCause.(*horizonclient.Error).Problem)
							return
						default:
							log.Println(err.Error())
							return
						}
					}
				}
				if err = db.AddBlock(*NewBlock(blockNumber)); err != nil {
					log.WithError(err).
						WithField("block_number", blockNumber).
						Warnf("failed to add block")
					return
				}

			}()
		}
	}()

	return NewRootHandler(actService)
}

const int64Len = 19

func truncateToStellarPrecision(number *big.Int) int {
	str := number.String()
	length := len(str)

	sliceLen := int(math.Min(int64Len, float64(length)))
	slice := str[:sliceLen]

	amount, _ := strconv.Atoi(slice)

	return amount
}
