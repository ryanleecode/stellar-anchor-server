package internal

import (
	"github.com/go-errors/errors"
	"os"
	"strconv"
)

type Environment struct {
	port              string
	mnemonic          string
	dbHost            string
	dbPort            string
	dbUser            string
	dbName            string
	dbPassword        string
	dbSSLMode         string
	ethIPCEndpoint    string
	networkPassphrase string
	amqpPort          string
	amqpHost          string
	amqpUser          string
	amqpPassword      string
}

func (e *Environment) Port() string {
	return e.port
}

func (e *Environment) Mnemonic() string {
	return e.mnemonic
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

func (e *Environment) EthIPCEndpoint() string {
	return e.ethIPCEndpoint
}

func (e *Environment) NetworkPassphrase() string {
	return e.networkPassphrase
}

func (e *Environment) AMQPPort() string {
	return e.amqpPort
}

func (e *Environment) AMQPHost() string {
	return e.amqpHost
}

func (e *Environment) AMQPUser() string {
	return e.amqpUser
}

func (e *Environment) AMQPPassword() string {
	return e.amqpPassword
}

func NewEnvironment() *Environment {
	return &Environment{
		port:              os.Getenv("PORT"),
		mnemonic:          os.Getenv("MNEMONIC"),
		dbHost:            os.Getenv("DB_HOST"),
		dbPort:            os.Getenv("DB_PORT"),
		dbUser:            os.Getenv("DB_USER"),
		dbName:            os.Getenv("DB_NAME"),
		dbPassword:        os.Getenv("DB_PASSWORD"),
		dbSSLMode:         os.Getenv("DB_SSL_MODE"),
		networkPassphrase: os.Getenv("NETWORK_PASSPHRASE"),
		ethIPCEndpoint:    os.Getenv("ETH_IPC_ENDPOINT"),
		amqpPort:          os.Getenv("AMQP_PORT"),
		amqpHost:          os.Getenv("AMQP_HOST"),
		amqpPassword:      os.Getenv("AMQP_PASSWORD"),
		amqpUser:          os.Getenv("AMQP_USER"),
	}
}

func (e *Environment) Validate() []error {
	var errs []error
	if _, err := strconv.ParseUint(e.port, 10, 32); err != nil {
		errs = append(errs, errors.New("PORT is not a valid number"))
	}
	if e.mnemonic == "" {
		errs = append(errs, errors.New("MNEMONIC is missing"))
	}
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
	if e.networkPassphrase == "" {
		errs = append(errs, errors.New("NETWORK_PASSPHRASE is missing"))
	}
	if e.ethIPCEndpoint == "" {
		errs = append(errs, errors.New("ETH_IPC_ENDPOINT is missing"))
	}
	if e.amqpPort == "" {
		errs = append(errs, errors.New("AMQP_PORT is missing"))
	}
	if e.amqpHost == "" {
		errs = append(errs, errors.New("AMQP_HOST is missing"))
	}
	if e.amqpUser == "" {
		errs = append(errs, errors.New("AMQP_USER is missing"))
	}
	if e.amqpPassword == "" {
		errs = append(errs, errors.New("AMQP_PASSWORD is missing"))
	}

	return errs
}
