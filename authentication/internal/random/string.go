package random

import (
	"encoding/base64"
	"github.com/pkg/errors"
)

type StringGenerator func(int) (string, error)

// NewGenerateString returns a function that provides
// URL-safe, base64 encoded securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func NewGenerateString(generateBytes func(n int) ([]byte, error)) StringGenerator {
	return func(s int) (string, error) {
		b, err := generateBytes(s)
		if err != nil {
			return "", errors.Wrap(err, "cannot generate random string")
		}

		return base64.URLEncoding.EncodeToString(b), nil
	}
}
