package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

func main() {
	http.Handle("/", &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = "127.0.0.1:9222"
		},
	})
	host, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:9222", host), nil))
}
