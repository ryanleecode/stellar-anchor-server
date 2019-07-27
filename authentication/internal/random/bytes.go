package random

import (
	"github.com/pkg/errors"
)

type ByteReader func([]byte) (int, error)
type ByteGenerator func(int) ([]byte, error)

func NewGenerateBytes(r ByteReader) ByteGenerator {
	return func(n int) (bytes []byte, e error) {
		b := make([]byte, n)
		_, err := r(b)
		// Note that err == nil only if we read len(b) bytes.
		if err != nil {
			return nil, errors.Wrap(err, "cannot generate random bytes")
		}

		return b, nil
	}
}
