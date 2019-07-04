package authentication

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"time"
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
	txe *xdr.TransactionEnvelope,
) []error {
	tx := txe.Tx
	txAnchorPubKey := tx.SourceAccount.Address()
	timebounds := tx.TimeBounds
	operations := tx.Operations
	validationErrs := make([]error, 0)
	if txAnchorPubKey != s.keypair.Address() {
		validationErrs = append(validationErrs, NewTransactionSourceAccountDoesntMatchAnchorPublicKey(
			fmt.Sprintf("the transaction's address does not match the anchor's")))
	}
	if timebounds == nil {
		validationErrs = append(
			validationErrs, NewTransactionIsMissingTimeBounds("transaction is missing timebounds"))
	} else {
		now := xdr.TimePoint(time.Now().UTC().Unix())
		if now > timebounds.MaxTime || now < timebounds.MinTime {
			validationErrs = append(validationErrs, NewTransactionChallengeExpired("transaction challenge has expired"))
		}
	}

	if operations == nil {
		validationErrs = append(validationErrs, NewTransactionOperationsIsNil("transaction is missing manage data operation"))
	} else if len(operations) != 1 {
		validationErrs = append(validationErrs, NewTransactionChallengeDoesNotHaveOnlyOneOperation(
			fmt.Sprintf("transaction can only have one operation but found %d", len(operations))))
	} else {
		operation := operations[0]
		if operation.Body.Type != xdr.OperationTypeManageData {
			validationErrs = append(validationErrs, NewTransactionChallengeIsNotAManageDataOperation(
				"expected transaction to have a manage data operation type"))
		}
		if operation.SourceAccount == nil {
			validationErrs = append(validationErrs, NewTransactionOperationSourceAccountIsEmpty(
				"transaction operation does not have a source account id"))
		}

	}

	return validationErrs
}
