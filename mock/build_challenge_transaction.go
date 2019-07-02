package mock

import (
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/mock"
)

type BuildChallengeTransactionMock struct {
	mock.Mock
}

func (m *BuildChallengeTransactionMock) BuildChallengeTransaction(
	serverAccount *horizon.Account, clientAccount *horizon.Account,
) (*txnbuild.Transaction, error) {
	args := m.Called(serverAccount, clientAccount)

	return args.Get(0).(*txnbuild.Transaction), args.Error(1)
}
