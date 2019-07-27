package internal

import (
	"github.com/pkg/errors"
	"os"
	"strconv"
)

type Environment struct {
	port string
	networkPassphrase string
	mnemonic string
}

func NewEnvironment() *Environment {
	return &Environment{
		port: os.Getenv("PORT"),
		networkPassphrase: os.Getenv("NETWORK_PASSPHRASE"),
		mnemonic: os.Getenv("MNEMONIC"),
	}
}

func (e *Environment) Port() string {
	return e.port
}

func (e *Environment) NetworkPassphrase() string {
	return e.networkPassphrase
}

func (e *Environment) Mnemonic() string {
	return e.mnemonic
}

func (e *Environment) Validate() []error {
	var errs []error
	if _, err := strconv.ParseUint(e.port, 10, 32); err != nil {
		errs = append(errs, errors.New("PORT is not a valid number"))
	}
	if e.networkPassphrase == "" {
		errs = append(errs, errors.New("NETWORK_PASSPHRASE key is missing"))
	}
	if e.mnemonic == "" {
		errs = append(errs, errors.New("MNEMONIC key is missing"))
	}

	return errs
}