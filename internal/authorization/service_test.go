package authorization

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"stellar-fi-anchor/mock"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
}

func (s *ServiceSuite) SetupTest() {

}

func (s *ServiceSuite) TestValidationFailsWhenSourceAccountDoesntMatchPublicKey() {
	stellarClientMock := new(mock.StellarClientMock)
	buildChallengeTransactionMock := new(mock.BuildChallengeTransactionMock)
	anchorKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)
	authService := NewService(
		stellarClientMock,
		buildChallengeTransactionMock.BuildChallengeTransaction,
		anchorKeyPair)

	randomKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)
	err = authService.ValidateClientSignedChallengeTransaction(
		randomKeyPair.Address())
	assert.Error(s.T(), err)
	origErr := errors.Cause(err)
	assert.IsType(s.T(), &TransactionSourceAccountDoesntMatchAnchorPublicKey{}, origErr)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
