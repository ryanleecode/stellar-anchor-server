package random

import (
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestByteGeneratorSuccess(t *testing.T) {
	rand.Seed(1337)
	generator := NewGenerateBytes(rand.Read)

	bytes, err := generator(6)
	assert.NoError(t, err)
	assert.EqualValues(t, bytes, []byte{0x26, 0xc5, 0xa4, 0x18, 0x2a, 0x81})
}

func TestByteGeneratorFailure(t *testing.T) {
	byteReader := func([]byte) (int, error) {
		return 0, errors.New("cannot read bytes")
	}
	generator := NewGenerateBytes(byteReader)

	bytes, err := generator(6)
	assert.Error(t, err)
	assert.Nil(t, bytes)
}
