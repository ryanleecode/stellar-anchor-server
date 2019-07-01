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
	"log"
	"net/http"
	"time"

	"stellar-fi-anchor/internal"
)

func main() {
	rootHandler := internal.NewRootHandler()

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %d\n", 8000)
	log.Fatal(server.ListenAndServe())
}
