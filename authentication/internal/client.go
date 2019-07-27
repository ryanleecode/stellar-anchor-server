package internal

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
)

type Client struct {
	client *horizonclient.Client
}

func NewClient(c *horizonclient.Client) *Client {
	return &Client{client: c}
}

func (c *Client) AccountDetail(request horizonclient.AccountRequest) (*horizon.Account, error) {
	account, err := c.client.AccountDetail(request)
	if err != nil {
		problem := err.(*horizonclient.Error).Problem
		if horizonBadRequestRegex.MatchString(problem.Type) {
			if problem.Extras["invalid_field"] == "account_id" {
				return nil, &InvalidAccountID{message: "invalid account id"}
			}
		}
		return nil, errors.Wrap(err, "failed to retrieve account details")
	}

	return &account, nil
}
