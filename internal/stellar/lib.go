package stellar

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/txnbuild"
	"time"
)

type Client struct {
	client *horizonclient.Client
}

func NewClient(c *horizonclient.Client) *Client {
	return &Client{client: c}
}

func BuildChallengeTransaction(
	serverAccount *horizon.Account, clientAccount *horizon.Account,
) (*txnbuild.Transaction, error) {
	currentSequenceNumber := serverAccount.Sequence
	defer func() {
		serverAccount.Sequence = currentSequenceNumber
	}()
	serverAccount.Sequence = "0"
	randomNounce, err := GenerateRandomString(48)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate random nounce for challenge")
	}

	currentTime := time.Now().UTC().Unix()
	tx := txnbuild.Transaction{
		SourceAccount: serverAccount,
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
		Network: network.TestNetworkPassphrase,
	}
	err = tx.Build()
	if err != nil {
		return nil, errors.Wrap(err, "cannot build challenge txn")
	}

	return &tx, nil
}
