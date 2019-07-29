package internal

import (
	"errors"
	"os"
	"strconv"
)

type Environment struct {
	port string
}

func NewEnvironment() *Environment {
	return &Environment{
		port: os.Getenv("PORT"),
	}
}

func (e *Environment) Port() uint {
	p, _ := strconv.ParseUint(e.port, 10, 32)

	return uint(p)
}

func (e *Environment) Validate() []error {
	var errs []error

	if _, err := strconv.ParseUint(e.port, 10, 32); err != nil {
		errs = append(errs, errors.New("PORT is not a valid number"))
	}

	return errs
}