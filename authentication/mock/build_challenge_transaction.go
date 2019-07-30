package mock

import (
	"github.com/drdgvhbh/stellar-anchor-server/authentication/internal"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/mock"
)

type BuildChallengeTransactionMock struct {
	mock.Mock
}

func (m *BuildChallengeTransactionMock) Build(
	serverAccount internal.Account, clientAccount internal.Account,
) (*txnbuild.Transaction, error) {
	args := m.Called(serverAccount, clientAccount)

	return args.Get(0).(*txnbuild.Transaction), args.Error(1)
}
