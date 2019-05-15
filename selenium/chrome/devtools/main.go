package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	listen string
)

func init() {
	flag.StringVar(&listen, "listen", ":7070", "Network address to accept connections")
	flag.Parse()
}

func main() {
	log.Printf("[INIT] [Listening on %s]", listen)
	log.Fatal(http.ListenAndServe(listen, ws()))
}
