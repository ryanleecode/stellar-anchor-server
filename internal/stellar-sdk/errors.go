package stellarsdk

import (
	"regexp"
)

var horizonBadRequestRegex = regexp.MustCompile("bad_request")

type InvalidAccountID struct {
	message string
}

func (e *InvalidAccountID) Error() string {
	return e.message
}
