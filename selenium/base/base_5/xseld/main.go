package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cmd := exec.Command("xsel", "-b")
		switch r.Method {
		case http.MethodGet:
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			go io.Copy(w, stdout)
		case http.MethodPost:
			cmd.Args = append(cmd.Args, "-i")
			stdin, err := cmd.StdinPipe()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			go func() {
				defer stdin.Close()
				io.Copy(stdin, r.Body)
			}()
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		err := cmd.Run()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}
