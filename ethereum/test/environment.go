package test

import (
	"errors"
	"os"
	"strconv"
)

type Environment struct {
	dbHost         string
	dbPort         string
	dbUser         string
	dbName         string
	dbPassword     string
	dbSSLMode      string
	ethRPCEndpoint string
}

func (e *Environment) DBHost() string {
	return e.dbHost
}

func (e *Environment) DBPort() string {
	return e.dbPort
}

func (e *Environment) DBUser() string {
	return e.dbUser
}

func (e *Environment) DBName() string {
	return e.dbName
}

func (e *Environment) DBPassword() string {
	return e.dbPassword
}

func (e *Environment) DBSSLMode() string {
	return e.dbSSLMode
}

func (e *Environment) EthRPCEndpoint() string {
	return e.ethRPCEndpoint
}

func NewEnvironment() *Environment {
	return &Environment{
		dbHost:         os.Getenv("DB_HOST"),
		dbPort:         os.Getenv("DB_PORT"),
		dbUser:         os.Getenv("DB_USER"),
		dbName:         os.Getenv("DB_NAME"),
		dbPassword:     os.Getenv("DB_PASSWORD"),
		dbSSLMode:      os.Getenv("DB_SSL_MODE"),
		ethRPCEndpoint: os.Getenv("ETH_RPC_ENDPOINT"),
	}
}

func (e *Environment) Validate() []error {
	var errs []error
	if e.dbHost == "" {
		errs = append(errs, errors.New("DB_HOST is missing"))
	}
	if _, err := strconv.ParseUint(e.dbPort, 10, 32); err != nil {
		errs = append(errs, errors.New("DB_PORT is not a valid number"))
	}
	if e.dbName == "" {
		errs = append(errs, errors.New("DB_NAME is missing"))
	}
	if e.dbUser == "" {
		errs = append(errs, errors.New("DB_USER is missing"))
	}
	if e.dbSSLMode != "disable" && e.dbSSLMode != "enable" {
		errs = append(errs, errors.New("DB_SSL_MODE must be disable or enable"))
	}
	if e.ethRPCEndpoint == "" {
		errs = append(errs, errors.New("ETH_RPC_ENDPOINT is missing"))
	}

	return errs
}
