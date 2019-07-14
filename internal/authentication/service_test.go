package authentication

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
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
	passphrase                    string
}

func (s *ServiceSuite) SetupTest() {
	s.stellarClientMock = new(mock.StellarClientMock)
	s.buildChallengeTransactionMock = new(mock.BuildChallengeTransactionMock)

	anchorKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)

	s.anchorKeyPair = anchorKeyPair
	s.passphrase = network.TestNetworkPassphrase

	s.authService = NewService(
		s.stellarClientMock,
		s.buildChallengeTransactionMock,
		s.anchorKeyPair,
		s.passphrase)
}

func (s *ServiceSuite) generateChallengeTransaction(
	clientAddress string,
	timebounds *txnbuild.Timebounds,
	operations []txnbuild.Operation,
) *txnbuild.Transaction {
	tx := txnbuild.Transaction{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: s.anchorKeyPair.Address(),
			Sequence:  -1,
		},
		Operations: operations,
		Network:    s.passphrase,
	}
	if timebounds != nil {
		tx.Timebounds = *timebounds
	} else {
		tx.Timebounds = txnbuild.NewInfiniteTimeout()
	}

	return &tx
}

func (s *ServiceSuite) TestValidationFailsWhenSourceAccountDoesntMatchPublicKey() {
	clientKP, err := keypair.Random()
	assert.NoError(s.T(), err)

	incorrectAnchorKP, err := keypair.Random()
	assert.NoError(s.T(), err)

	tx := s.generateChallengeTransaction(clientKP.Address(), nil, nil)
	tx.SourceAccount = &txnbuild.SimpleAccount{
		AccountID: incorrectAnchorKP.Address(),
	}
	err = tx.Build()
	assert.NoError(s.T(), err)

	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
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
	err := tx.Build()
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()
	txEnv.Tx.TimeBounds = nil

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
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
	now := time.Now().UTC().Unix()
	timeBounds := txnbuild.NewTimebounds(now-3, now-1)

	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), &timeBounds, nil)
	err := tx.Build()
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
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
	now := time.Now().UTC().Unix()
	timeBounds := txnbuild.NewTimebounds(now+1, now+3)

	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), &timeBounds, nil)
	err := tx.Build()
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
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
		s.anchorKeyPair.Address(), nil, []txnbuild.Operation{
			&txnbuild.BumpSequence{},
			&txnbuild.BumpSequence{},
		})
	err := tx.Build()
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
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
	ops := []txnbuild.Operation{
		&txnbuild.BumpSequence{},
	}
	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), nil, ops)
	err := tx.Build()
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
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
	err := tx.Build()
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
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
	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), nil, nil)
	err := tx.Build()
	assert.NoError(s.T(), err)
	randomKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)
	err = tx.Sign(randomKeyPair)
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionIsNotSignedByAnchor:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsIfTransactionIsSignedByAnchorButWithTheWrongPassphrase() {
	s.passphrase = "private network 12345 - no haxor plz"
	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), nil, nil)
	err := tx.Build()
	assert.NoError(s.T(), err)
	randomKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)
	err = tx.Sign(randomKeyPair)
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionIsNotSignedByAnchor:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsIfTransactionIsNotSignedByClient() {
	clientKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)

	ops := []txnbuild.Operation{
		&txnbuild.ManageData{
			SourceAccount: &txnbuild.SimpleAccount{AccountID: clientKeyPair.Address()},
		},
	}

	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), nil, ops)
	err = tx.Build()
	assert.NoError(s.T(), err)
	randomKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)
	err = tx.Sign(randomKeyPair)
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionIsNotSignedByClient:
			return true
		default:
			return false
		}
	})
	assert.True(s.T(),
		len(filteredErrs.([]error)) == 1)
}

func (s *ServiceSuite) TestValidationFailsIfTransactionIsByClientButWithTheWrongPassphrase() {
	s.passphrase = "private network 12345 - no haxor plz"
	clientKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)

	ops := []txnbuild.Operation{
		&txnbuild.ManageData{
			SourceAccount: &txnbuild.SimpleAccount{AccountID: clientKeyPair.Address()},
		},
	}

	tx := s.generateChallengeTransaction(
		s.anchorKeyPair.Address(), nil, ops)
	err = tx.Build()
	assert.NoError(s.T(), err)
	randomKeyPair, err := keypair.Random()
	assert.NoError(s.T(), err)
	err = tx.Sign(randomKeyPair)
	assert.NoError(s.T(), err)
	txEnv := tx.TxEnvelope()

	validationErrs := s.authService.ValidateClientSignedChallengeTransaction(txEnv)
	filteredErrs := funk.Filter(validationErrs, func(x error) bool {
		origErr := errors.Cause(x)
		switch origErr.(type) {
		case *TransactionIsNotSignedByClient:
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
