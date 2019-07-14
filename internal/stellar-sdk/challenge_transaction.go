package stellarsdk

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/txnbuild"
	"time"
)

type ChallengeTransactionFactory struct {
	passphrase      string
	nounceGenerator func() (string, error)
}

func NewChallengeTransactionFactory(networkPassphrase string, nounceGenerator func() (string, error)) *ChallengeTransactionFactory {
	return &ChallengeTransactionFactory{
		passphrase:      networkPassphrase,
		nounceGenerator: nounceGenerator,
	}
}

type Account txnbuild.Account

func (f *ChallengeTransactionFactory) Build(serverAccount Account, clientAccount Account) (*txnbuild.Transaction, error) {
	randomNounce, err := f.nounceGenerator()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate random nounce for challenge")
	}

	currentTime := time.Now().UTC().Unix()
	tx := txnbuild.Transaction{
		SourceAccount: &txnbuild.SimpleAccount{
			AccountID: serverAccount.GetAccountID(),
			Sequence:  -1,
		},
		Timebounds: txnbuild.NewTimebounds(
			currentTime,
			currentTime+(int64(time.Minute.Seconds())*5),
		),
		Operations: []txnbuild.Operation{
			&txnbuild.ManageData{
				SourceAccount: clientAccount,
				Name:          "Stellar FI Anchor auth",
				Value:         []byte(randomNounce),
			},
		},
		Network: f.passphrase,
	}
	err = tx.Build()
	if err != nil {
		return nil, errors.Wrap(err, "cannot build challenge txn")
	}

	return &tx, nil
}
