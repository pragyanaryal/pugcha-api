package router

import (
	"github.com/gorilla/mux"
	"gitlab.com/ProtectIdentity/pugcha-backend/middleware"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func SetupFrontRoutes(r *mux.Router, host string) {
	front := r.Host(host).Subrouter()

	target := "http://127.0.0.1:3000"
	remote, err := url.Parse(target)
	if err != nil {
		panic(err)
	}


	proxy := httputil.NewSingleHostReverseProxy(remote)
	front.HandleFunc("/{rest:.*}", handler(proxy, remote))

	front.Use(middleware.PanicRecover)
	front.Use(middleware.LogEachRequest)
}


func handler(p *httputil.ReverseProxy, url2 *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Host = url2.Host
		r.URL.Scheme = url2.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = url2.Host
		r.URL.Path = mux.Vars(r)["rest"]
		p.ServeHTTP(w, r)
	}
}