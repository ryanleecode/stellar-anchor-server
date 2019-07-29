package internal

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRootHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/.well-known/stellar.toml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "stellar.toml")
	})

	return handlers.RecoveryHandler()(router)
}