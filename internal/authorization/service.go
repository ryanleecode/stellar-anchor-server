package authorization

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

type StellarClient interface {
	AccountDetail(request horizonclient.AccountRequest) (*horizon.Account, error)
}

type BuildChallengeTransaction func(serverAccount *horizon.Account, clientAccount *horizon.Account) (*txnbuild.Transaction, error)

type Service struct {
	stellarClient StellarClient
	build         BuildChallengeTransaction
	keypair       *keypair.Full
}

func NewService(stellarClient StellarClient, build BuildChallengeTransaction, keypair *keypair.Full) *Service {
	return &Service{
		stellarClient: stellarClient,
		build:         build,
		keypair:       keypair,
	}
}

func (s *Service) BuildSignEncodeChallengeTransactionForAccount(id string) (string, error) {
	clientAccountRequest := horizonclient.AccountRequest{AccountID: id}
	clientAccount, err := s.stellarClient.AccountDetail(clientAccountRequest)
	if err != nil {
		return "", errors.Wrap(err, "cannot fetch client account details")
	}
	anchorPublicKey := s.keypair.Address()
	serverAccountRequest := horizonclient.AccountRequest{AccountID: anchorPublicKey}
	serverAccount, err := s.stellarClient.AccountDetail(serverAccountRequest)
	if err != nil {
		return "", errors.Wrap(err, "cannot fetch server account details")
	}
	txn, err := s.build(serverAccount, clientAccount)
	if err != nil {
		return "", errors.Wrap(err, "cannot build challenge txn")
	}
	err = txn.Sign(s.keypair)
	if err != nil {
		return "", errors.Wrap(err, "cannot sign challenge txn")
	}
	b64e, err := txn.Base64()
	if err != nil {
		return "", errors.Wrap(err, "cannot base64 encode signed challenge txn")
	}

	return b64e, nil
}

func (s *Service) ValidateClientSignedChallengeTransaction(
	anchorPublicKey string,
	timebounds *xdr.TimeBounds,
) error {
	if anchorPublicKey != s.keypair.Address() {
		return NewTransactionSourceAccountDoesntMatchAnchorPublicKey(
			fmt.Sprintf("the transaction's address does not match the anchor's"))
	}
	if timebounds == nil {
		return NewTransactionIsMissingTimeBounds("transaction is missing timebounds")
	}
	return nil
}
