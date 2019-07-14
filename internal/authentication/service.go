package authentication

import (
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	stellarsdk "stellar-fi-anchor/internal/stellar-sdk"
	"time"
)

type ChallengeTransactionFactory interface {
	Build(serverAccount stellarsdk.Account, clientAccount stellarsdk.Account) (*txnbuild.Transaction, error)
}

type Service struct {
	challengeTxFactory ChallengeTransactionFactory
	keypair            *keypair.Full
	networkPassphrase  string
}

func NewService(fact ChallengeTransactionFactory, kp *keypair.Full, passphrase string) *Service {
	return &Service{
		challengeTxFactory: fact,
		keypair:            kp,
		networkPassphrase:  passphrase,
	}
}

func (s *Service) BuildSignEncodeChallengeTransactionForAccount(id string) (string, error) {
	clientAccount := txnbuild.SimpleAccount{
		AccountID: id,
		Sequence:  -1,
	}

	anchorPublicKey := s.keypair.Address()
	serverAccount := txnbuild.SimpleAccount{
		AccountID: anchorPublicKey,
		Sequence:  -1,
	}
	txn, err := s.challengeTxFactory.Build(&serverAccount, &clientAccount)

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
	txAnchorPK := tx.SourceAccount.Address()
	timebounds := tx.TimeBounds
	operations := tx.Operations
	validationErrs := make([]error, 0)
	anchorPK := s.keypair.Address()
	if txAnchorPK != anchorPK {
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

	clientPublicKey := ""
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
		} else {
			clientPublicKey = operation.SourceAccount.Address()
		}
	}

	hash, err := network.HashTransaction(&tx, s.networkPassphrase)
	if err != nil {
		validationErrs = append(validationErrs, errors.Wrap(err, "cannot hash transaction"))
		return validationErrs
	}

	isSignedByAnchor := validateTransactionIsSignedBy(s.keypair, hash[:], txe.Signatures)
	if !isSignedByAnchor {
		validationErrs = append(validationErrs, NewTransactionIsNotSignedByAnchor(
			"transaction is not signed by the anchor"))
	}

	clientKeyPair, err := keypair.Parse(clientPublicKey)
	if err != nil {
		validationErrs = append(validationErrs, NewCannotParseClientPublicKey(
			fmt.Sprintf("cannot parse client public key %s", clientPublicKey)))
	}
	if clientKeyPair != nil {
		isSignedByClient := validateTransactionIsSignedBy(clientKeyPair, hash[:], txe.Signatures)
		if !isSignedByClient {
			validationErrs = append(validationErrs, NewTransactionIsNotSignedByClient(
				"transaction is not signed by the client"))
		}

	}

	return validationErrs
}

func validateTransactionIsSignedBy(
	kp keypair.KP,
	transaction []byte,
	signatures []xdr.DecoratedSignature,
) bool {
	if transaction == nil || signatures == nil || len(signatures) == 0 {
		return false
	}

	for _, decorSig := range signatures {
		err := kp.Verify(transaction, decorSig.Signature)
		if err == nil {
			return true
		}
	}

	return false
}

func (s *Service) Authenticate(txe *xdr.TransactionEnvelope) (string, error) {
	now := time.Now()
	txeHash, err := network.HashTransaction(&txe.Tx, s.networkPassphrase)
	if err != nil {
		return "", errors.Wrap(err, "transaction cannot be hashed")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "", // TODO: Assign iss from pulling stellar.toml
		"sub": txe.Tx.Operations[0].SourceAccount.Address(),
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
		"jti": hex.EncodeToString(txeHash[:]),
	})
	signedJwt, err := token.SignedString([]byte(s.keypair.Seed()))
	if err != nil {
		return "", errors.Wrap(err, "failed to sign jwt token")
	}

	return signedJwt, nil
}
