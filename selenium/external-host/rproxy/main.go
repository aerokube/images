package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	u, err := url.Parse(os.Args[1])
	if err != nil {
		log.Fatal("parse url: %v", err)
	}
	log.Fatal(http.ListenAndServe(":4444", httputil.NewSingleHostReverseProxy(u)))
}
