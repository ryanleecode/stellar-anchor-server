package authorization

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	err = authService.ValidateClientSignedChallengeTransaction(
		randomKeyPair.Address(),
		nil)
	assert.Error(s.T(), err)
	origErr := errors.Cause(err)
	assert.IsType(s.T(), &TransactionSourceAccountDoesntMatchAnchorPublicKey{}, origErr)
}

func (s *ServiceSuite) TestValidationFailsWhenTimeboundsIsNil() {
	authService := NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock.BuildChallengeTransaction,
		s.anchorKeyPair)

	err := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		nil)
	assert.Error(s.T(), err)
	origErr := errors.Cause(err)
	assert.IsType(s.T(), &TransactionIsMissingTimeBounds{}, origErr)
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
	err := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		&timeBounds)
	assert.Error(s.T(), err)
	origErr := errors.Cause(err)
	assert.IsType(s.T(), &TransactionChallengeExpired{}, origErr)
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
	err := authService.ValidateClientSignedChallengeTransaction(
		s.anchorKeyPair.Address(),
		&timeBounds)
	assert.Error(s.T(), err)
	origErr := errors.Cause(err)
	assert.IsType(s.T(), &TransactionChallengeExpired{}, origErr)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
