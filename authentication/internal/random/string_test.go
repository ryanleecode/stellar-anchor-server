package random

import (
	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockByteGenerator struct {
	mock.Mock
}

func (g *MockByteGenerator) Generate(n int) ([]byte, error) {
	args := g.Called(n)

	return args.Get(0).([]byte), args.Error(1)
}

func TestStringGeneratorSuccess(t *testing.T) {
	byteGenerator := new(MockByteGenerator)
	numBytes := 4
	byteGenerator.On("Generate", numBytes).Return([]byte{1, 2, 3, 4}, nil)

	bytes, err := NewGenerateString(byteGenerator.Generate)(numBytes)
	assert.Equal(t, "AQIDBA==", bytes)
	assert.NoError(t, err)
}

func TestStringGeneratorError(t *testing.T) {
	byteGenerator := new(MockByteGenerator)
	numBytes := 4
	byteGenerator.On("Generate", numBytes).Return(
		[]byte{}, errors.New("byte generator failed"))

	bytes, err := NewGenerateString(byteGenerator.Generate)(numBytes)
	assert.Equal(t, "", bytes)
	assert.Error(t, err)
}
