package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	listen string
	android bool
)

func init() {
	flag.StringVar(&listen, "listen", ":7070", "Network address to accept connections")
	flag.BoolVar(&android, "android", false, "Whether we need to forward Android Emulator devtools port")
}

func main() {
	flag.Parse()
	log.Printf("[INIT] [Listening on %s]", listen)
	log.Fatal(http.ListenAndServe(listen, root()))
}
