package internal

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type APIGatewayParams interface {
	AuthServer() *url.URL
	StaticServer() *url.URL
}

func NewRootHandler(params APIGatewayParams) http.Handler {
	root := mux.NewRouter()
	root.HandleFunc("/.well-known/stellar.toml",
		proxy("", httputil.NewSingleHostReverseProxy(params.StaticServer())))

	apiRouter := root.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/authentications",
		proxy("authentications", httputil.NewSingleHostReverseProxy(params.AuthServer())))

	return root
}

func proxy(path string, p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		if path != "" {
			r.URL.Path = path
		}

		p.ServeHTTP(w, r)
	}
}