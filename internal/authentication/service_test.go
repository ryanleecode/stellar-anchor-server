package authentication

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thoas/go-funk"
	"stellar-fi-anchor/mock"
	"testing"
	"time"
)

type ServiceSuite struct {
	suite.Suite
	stellarClientMock             *mock.StellarClientMock
	buildChallengeTransactionMock *mock.BuildChallengeTransactionMock
	anchorKeyPair                 *keypair.Full
	authService                   *Service
}

func (s *ServiceSuite) SetupTest() {
	s.stellarClientMock = new(mock.StellarClientMock)
	s.buildChallengeTransactionMock = new(mock.BuildChallengeTransactionMock)

	anchorKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)

	s.anchorKeyPair = anchorKeyPair

	s.authService = NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)
}

func (s *ServiceSuite) generateChallengeTransaction(
	clientAddress string,
	timebounds *xdr.TimeBounds,
	operations []xdr.Operation,
) *xdr.Transaction {
	pubKey := [32]byte{}
	copy(pubKey[:], s.anchorKeyPair.Address())

	sourceAccountKey, err := xdr.NewPublicKey(xdr.PublicKeyTypePublicKeyTypeEd25519, xdr.Uint256(pubKey))
	assert.NoError(s.T(), err)

	tx := xdr.Transaction{
		SourceAccount: xdr.AccountId(sourceAccountKey),
		TimeBounds:    timebounds,
		Operations:    operations,
	}

	return &tx
}

func (s *ServiceSuite) TestValidationFailsWhenSourceAccountDoesntMatchPublicKey() {
	randomKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)

	tx := s.generateChallengeTransaction(randomKeyPair.Address(), nil, nil)
	txEnv := xdr.TransactionEnvelope{
		Tx: *tx,
	}

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(&txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionSourceAccountDoesntMatchAnchorPublicKey:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsWhenTimeboundsIsNil() {
	tx := s.generateChallengeTransaction(s.anchorKeyPair.Address(), nil, nil)
	txEnv := xdr.TransactionEnvelope{
		Tx: *tx,
	}

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(&txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionIsMissingTimeBounds:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsWhenNowIsAfterTimeboundsMaxTime() {
	now := xdr.TimePoint(time.Now().UTC().Unix())
	timeBounds := xdr.TimeBounds{
		MinTime: now - 3,
		MaxTime: now - 1,
	}

	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), &timeBounds, nil)
	txEnv := xdr.TransactionEnvelope{
		Tx: *tx,
	}

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(&txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionChallengeExpired:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsWhenNowIsBeforeTimeboundsMinTime() {
	now := xdr.TimePoint(time.Now().UTC().Unix())
	timeBounds := xdr.TimeBounds{
		MinTime: now + 1,
		MaxTime: now + 3,
	}

	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), &timeBounds, nil)
	txEnv := xdr.TransactionEnvelope{
		Tx: *tx,
	}

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(&txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionChallengeExpired:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsIfThereIsNotOnlyOneOperation() {
	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), nil, []xdr.Operation{})
	txEnv := xdr.TransactionEnvelope{
		Tx: *tx,
	}

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(&txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionChallengeDoesNotHaveOnlyOneOperation:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsIfOperationIsNotAManageDataOperation() {
	ops := []xdr.Operation{{Body: xdr.OperationBody{Type: xdr.OperationTypePayment}}}
	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), nil, ops)
	txEnv := xdr.TransactionEnvelope{
		Tx: *tx,
	}

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(&txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionChallengeIsNotAManageDataOperation:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsIfOperationSourceAccountIsNil() {
	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), nil, nil)
	txEnv := xdr.TransactionEnvelope{
		Tx: *tx,
	}

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(&txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionOperationsIsNil:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsIfTransactionIsNotSignedByAnchor() {

}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
