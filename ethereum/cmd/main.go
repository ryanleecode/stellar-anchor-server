package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/stellar/go/network"

	"github.com/drdgvhbh/stellar-anchor-server/ethereum/internal"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger := log.New()
	logger.SetLevel(log.TraceLevel)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	env := internal.NewEnvironment()
	envErrors := env.Validate()
	if len(envErrors) > 0 {
		err := errors.New("")
		for _, e := range envErrors {
			err = errors.Wrapf(err, e.Error())
		}

		logger.Fatalln(err)
	}
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
		log.Fatalln(errors.Wrap(err, "failed to open database"))
	}
	ipcClient, err := rpc.DialIPC(context.Background(), env.EthIPCEndpoint())
	if err != nil {
		log.Fatalln(errors.Wrap(err, "failed to connect to ethereum ipc client"))
	}
	defer func() {
		_ = db.Close()
		ipcClient.Close()
	}()

	rootHandler := internal.Bootstrap(network.TestNetworkPassphrase, env.Mnemonic(), db, ipcClient, logger)

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", env.Port()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Printf("Server is listening on port %s", env.Port())
	logger.Fatal(server.ListenAndServe())
}
