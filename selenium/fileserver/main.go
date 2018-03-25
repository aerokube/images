package main

import (
	"flag"
	"log"
	"net/http"
)

var dir string

func init() {
	flag.StringVar(&dir, "dir", ".", "Directory to host")
	flag.Parse()
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir(dir))))
}
