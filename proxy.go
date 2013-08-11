package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type application struct {
	Name string
	Port string
	Test string
}

type proxy struct {
	apps      map[string]*application
	Address   string
	Transport *http.Transport
}

func NewProxy(address string) *proxy {
	p := &proxy{Address: address}
	p.Init()
	return p
}

func (p *proxy) Init() {
	p.apps = make(map[string]*application)

	mux := http.NewServeMux()
	p.Transport = &http.Transport{DisableKeepAlives: false, DisableCompression: false}

	mux.Handle("/", p)

	go func() {
		log.Printf("Starting proxy at %s\n", p.Address)
		srv := &http.Server{Handler: mux, Addr: p.Address}

		log.Fatal(srv.ListenAndServe())
	}()
}

func (p *proxy) Route(app *application) {
	p.apps[app.Name] = app

	log.Printf("Routing application `%s` to `%s`", app.Name, app.Port)
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")[1]
	app, ok := p.apps[path]

	if ok {
		r.URL.Scheme = "http"
		r.URL.Host = "localhost:" + app.Port

		resp, err := p.Transport.RoundTrip(r)

		if err != nil {
			p.responseError(err, w)
		} else {
			for k, v := range resp.Header {
				for _, vv := range v {
					w.Header().Add(k, vv)
				}
			}

			w.WriteHeader(resp.StatusCode)

			io.Copy(w, resp.Body)
			resp.Body.Close()
		}
	} else {
		p.responseError(fmt.Errorf("Not found"), w)
	}
}

func (p *proxy) responseError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintf(w, "Error: %v", err)
	log.Printf("Error: %v", err)
}
