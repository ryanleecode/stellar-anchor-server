package mock

import (
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/mock"
	stellarsdk "stellar-fi-anchor/internal/stellar-sdk"
)

type BuildChallengeTransactionMock struct {
	mock.Mock
}

func (m *BuildChallengeTransactionMock) Build(
	serverAccount stellarsdk.Account, clientAccount stellarsdk.Account,
) (*txnbuild.Transaction, error) {
	args := m.Called(serverAccount, clientAccount)

	return args.Get(0).(*txnbuild.Transaction), args.Error(1)
}
