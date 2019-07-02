package stellar

import (
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/support/errors"
	"regexp"
)

var horizonBadRequestRegex = regexp.MustCompile("bad_request")

type InvalidAccountID struct {
	message string
}

func (e *InvalidAccountID) Error() string {
	return e.message
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
