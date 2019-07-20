package main

import (
	"context"
	"fmt"
	"github.com/drdgvhbh/stellar-fi-anchor/internal"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

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
	db, err := gorm.Open(
		"postgres", fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
			environment.DBHost,
			environment.DBPort,
			environment.DBUser,
			environment.DBName,
			environment.DBSSLMode,
			environment.DBPassword))
	if err != nil {
		log.Fatalln(err, "failed to open database")
	}
	defer func() {
		_ = db.Close()
	}()

	client, err := rpc.DialIPC(context.Background(), "/home/drd/.ethereum/goerli/geth.ipc")
	if err != nil {
		log.Fatal(err)
	}

	rootHandler := internal.Bootstrap(environment.PrivateKey, environment.Mnemonic, db, client)

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", environment.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %s", environment.Port)
	log.Fatal(server.ListenAndServe())
}
