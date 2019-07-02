package mock

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stretchr/testify/mock"
)

type StellarClientMock struct {
	mock.Mock
}

func (m *StellarClientMock) AccountDetail(request horizonclient.AccountRequest) (*horizon.Account, error) {
	args := m.Called(request)

	return args.Get(0).(*horizon.Account), args.Error(1)
}
