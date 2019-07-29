package main

import (
	"fmt"
	"github.com/drdgvhbh/stellar-fi-anchor/authentication/internal"
	"github.com/pkg/errors"
	"log"
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

	rootHandler := internal.Bootstrap(environment)

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", environment.Port()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %s", environment.Port())
	log.Fatal(server.ListenAndServe())
}