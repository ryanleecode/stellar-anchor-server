package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"stellar-fi-anchor/internal"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Fatalln("env variable PORT not defined")
	}
	privateKey, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		log.Fatalln("env variable PRIVATE_KEY not defined")
	}
	mnemonic, ok := os.LookupEnv("MNEMONIC")
	if !ok {
		log.Fatalln("env variable MNEMONIC not defined")
	}
	db, err := gorm.Open(
		"postgres", "host=localhost port=6666 user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err, "failed to open database")
	}
	defer func() {
		_ = db.Close()
	}()

	rootHandler := internal.Bootstrap(privateKey, mnemonic, db)

	server := &http.Server{
		Handler:      rootHandler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Server is listening on port %d", 8000)
	log.Fatal(server.ListenAndServe())
}
