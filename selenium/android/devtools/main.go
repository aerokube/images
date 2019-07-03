package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
)

var (
	devtoolsUrl *url.URL
)

const (
	devtoolsPort = 9333
	listen = ":7070"
)

func init() {
	u, err := url.Parse(fmt.Sprintf("http://localhost:%d", devtoolsPort))
	if err != nil {
		log.Fatalf("[INIT] [Invalid devtools URL: %v]", err)
	}
	devtoolsUrl = u
}

func main() {
	log.Printf("[INIT] [Listening on %s]", listen)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("adb", "forward",  fmt.Sprintf("tcp:%d", devtoolsPort), "localabstract:chrome_devtools_remote")
		err := cmd.Run()
		if err != nil {
			log.Printf("[PORT_FORWARDING_ERROR] [%v]", err)
			http.Error(w, fmt.Sprintf("Failed to forward devtools port: %v", err), http.StatusInternalServerError)
			return
		}
		(&httputil.ReverseProxy{
			Director: func(r *http.Request) {
				log.Printf("[PROXY] [%s]", r.URL.String())
				r.URL = devtoolsUrl
			},
		}).ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(listen, nil))
}
