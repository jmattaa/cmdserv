package main

import (
	"log"
	"net/http"

	"github.com/jmattaa/cmdserv/middleware"
	"github.com/jmattaa/cmdserv/endpoints"
)

const DEBUG = true

func main() {
	if err := endpoints.Init(); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Handle(w, r)
	})

	var h http.Handler = mux
	if DEBUG {
		h = middleware.Logger(h)
	}

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", h))
}
