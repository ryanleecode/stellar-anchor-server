package internal

import (
	"github.com/pkg/errors"
	"net/url"
	"os"
	"strconv"
)

type Environment struct {
	port string
	authServer string
	staticServer string
}

func NewEnvironment() *Environment {
	return &Environment{
		port: os.Getenv("PORT"),
		authServer: os.Getenv("AUTH_SERVER"),
		staticServer: os.Getenv("STATIC_SERVER"),
	}
}

func (e *Environment) StaticServer() *url.URL {
	s, _ := url.Parse(e.staticServer)

	return s
}

func (e *Environment) AuthServer() *url.URL {
	s, _ := url.Parse(e.authServer)

	return s
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
	if _, err := url.Parse(e.authServer); err != nil {
		errs = append(errs, errors.New("AUTH_SERVER is not a valid url"))
	}
	if _, err := url.Parse(e.staticServer); err != nil {
		errs = append(errs, errors.New("STATIC_SERVER is not a valid url"))
	}

	return errs
}