// Package stellaranchor Stellar FI Anchor Emulator API.
//
// FI Anchor Emulator for the Stellar Network
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /v1
//     Version: 0.0.2
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Ryan Lee<ryanleecode@gmail.com> http://drdgvhbh.io
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"stellar-fi-anchor/internal"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalln("env variable PORT not defined")
	}
	config := internal.Config{
		Port: port,
	}

	rootHandler := internal.NewRootHandler()

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", config.GetPort()),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %d", 8000)
	log.Fatal(server.ListenAndServe())
}
