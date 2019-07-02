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
}

func (s *ServiceSuite) SetupTest() {
	s.stellarClientMock = new(mock.StellarClientMock)
	s.buildChallengeTransactionMock = new(mock.BuildChallengeTransactionMock)

	anchorKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)

	s.anchorKeyPair = anchorKeyPair
}

func (s *ServiceSuite) TestValidationFailsWhenSourceAccountDoesntMatchPublicKey() {
	authService := NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)

	randomKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)
	validationErrs := authService.ValidateClientSignedChallengeTransaction(
		randomKeyPair.Address(),
		nil,
		[]xdr.Operation{})
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
	authService := NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)

	validationErrs := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		nil,
		[]xdr.Operation{})
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
	authService := NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)

	now := xdr.TimePoint(time.Now().UTC().Unix())
	timeBounds := xdr.TimeBounds{
		MinTime: now - 3,
		MaxTime: now - 1,
	}
	validationErrs := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		&timeBounds,
		[]xdr.Operation{})
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
	authService := NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)

	now := xdr.TimePoint(time.Now().UTC().Unix())
	timeBounds := xdr.TimeBounds{
		MinTime: now + 1,
		MaxTime: now + 3,
	}
	validationErrs := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		&timeBounds,
		[]xdr.Operation{})
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
	authService := NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)

	now := xdr.TimePoint(time.Now().UTC().Unix())
	timeBounds := xdr.TimeBounds{
		MinTime: now - 1,
		MaxTime: now + 1,
	}
	validationErrs := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		&timeBounds,
		[]xdr.Operation{})
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
	authService := NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)

	now := xdr.TimePoint(time.Now().UTC().Unix())
	timeBounds := xdr.TimeBounds{
		MinTime: now - 1,
		MaxTime: now + 1,
	}
	validationErrs := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		&timeBounds,
		[]xdr.Operation{{Body: xdr.OperationBody{Type: xdr.OperationTypePayment}}})
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
	authService := NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)

	now := xdr.TimePoint(time.Now().UTC().Unix())
	timeBounds := xdr.TimeBounds{
		MinTime: now - 1,
		MaxTime: now + 1,
	}
	validationErrs := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		&timeBounds,
		[]xdr.Operation{{}})
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionOperationSourceAccountIsEmpty:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
