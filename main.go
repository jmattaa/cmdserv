package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/jmattaa/cmdserv/middleware"
)

const DEBUG = true

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := exec.Command("open", "/Applications/Spotify.app").Run()
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		}
	})

	var h http.Handler = mux
	if DEBUG {
		h = middleware.Logger(h)
	}

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", h))
}
