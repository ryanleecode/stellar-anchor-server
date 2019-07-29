package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/drdgvhbh/stellar-fi-anchor/ethereum/internal"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

type BootstrapParams struct {
	networkPassphrase string
	mnemonic          string
	db                *gorm.DB
	rpcClient         *rpc.Client
}

func NewBootstrapParams(env internal.Environment) (*BootstrapParams, error) {
	db, err := gorm.Open(
		"postgres", fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			env.DBHost(),
			env.DBPort(),
			env.DBUser(),
			env.DBName(),
			env.DBSSLMode(),
			env.DBPassword()))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	ipcClient, err := rpc.DialIPC(context.Background(), env.EthIPCEndpoint())
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to ethereum ipc client")
	}
	return &BootstrapParams{
		networkPassphrase: env.NetworkPassphrase(),
		mnemonic:          env.Mnemonic(),
		db:                db,
		rpcClient:         ipcClient,
	}, nil
}

func (p *BootstrapParams) NetworkPassphrase() string {
	return p.networkPassphrase
}

func (p *BootstrapParams) Mnemonic() string {
	return p.mnemonic
}

func (p *BootstrapParams) DB() *gorm.DB {
	return p.db
}

func (p *BootstrapParams) RPCClient() *rpc.Client {
	return p.rpcClient
}

func main() {
	environment := internal.NewEnvironment()
	envErrors := environment.Validate()
	if len(envErrors) > 0 {
		err := errors.New("")
		for _, e := range envErrors {
			err = errors.Wrapf(err, e.Error())
		}

		log.Fatalln(err)
	}
	bootstrapParams, err := NewBootstrapParams(*environment)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = bootstrapParams.DB().Close()
		bootstrapParams.RPCClient().Close()
	}()

	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	rootHandler := internal.Bootstrap(bootstrapParams)

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", environment.Port()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %s", environment.Port())
	log.Fatal(server.ListenAndServe())
}
