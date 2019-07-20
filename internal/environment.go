package internal

import (
	"github.com/go-errors/errors"
	"os"
	"strconv"
)

type Environment struct {
	Port           string
	PrivateKey     string
	Mnemonic       string
	DBHost         string
	DBPort         string
	DBUser         string
	DBName         string
	DBPassword     string
	DBSSLMode      string
	ETHIPCEndpoint string
}

func NewEnvironment() *Environment {
	return &Environment{
		Port:       os.Getenv("PORT"),
		PrivateKey: os.Getenv("PRIVATE_KEY"),
		Mnemonic:   os.Getenv("MNEMONIC"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBSSLMode:  os.Getenv("DB_SSL_MODE"),
	}
}

func (e *Environment) Validate() []error {
	var errs []error
	if _, err := strconv.ParseUint(e.Port, 10, 32); err != nil {
		errs = append(errs, errors.New("PORT is not a valid number"))
	}
	if e.PrivateKey == "" {
		errs = append(errs, errors.New("PRIVATE_KEY key is missing"))
	}
	if e.Mnemonic == "" {
		errs = append(errs, errors.New("MNEMONIC is missing"))
	}
	if e.DBHost == "" {
		errs = append(errs, errors.New("DB_HOST is missing"))
	}
	if _, err := strconv.ParseUint(e.DBPort, 10, 32); err != nil {
		errs = append(errs, errors.New("DB_PORT is not a valid number"))
	}
	if e.DBName == "" {
		errs = append(errs, errors.New("DB_NAME is missing"))
	}
	if e.DBUser == "" {
		errs = append(errs, errors.New("DB_USER is missing"))
	}
	if e.DBSSLMode != "disable" && e.DBSSLMode != "enable" {
		errs = append(errs, errors.New("DB_SSL_MODE must be disable or enable"))
	}

	return errs
}
