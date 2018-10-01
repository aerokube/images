package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", mux("/static")))
}

func mux(dir string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		if r.Method == http.MethodDelete {
			deleteFileIfExists(w, r, dir)
			return
		}
		http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
	})
	return mux
}

func deleteFileIfExists(w http.ResponseWriter, r *http.Request, dir string) {
	fileName := strings.TrimPrefix(r.URL.Path, "/")
	filePath := filepath.Join(dir, fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unknown file %s", fileName), http.StatusNotFound)
		return
	}
	err = os.Remove(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete file %s: %v", fileName, err), http.StatusInternalServerError)
		return
	}
}