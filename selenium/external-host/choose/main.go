package main

import (
	"os"
	"time"
	"fmt"
	"encoding/json"
	"math/rand"
)

var (
	s rand.Source
	r *rand.Rand
)

func init() {
	s = rand.NewSource(time.Now().Unix())
	r = rand.New(s)
}

func main() {
	var hosts []string
	if err := json.Unmarshal([]byte(os.Getenv("URLS")), &hosts); err != nil {
		fmt.Printf("unmarshall json: %v", err)
		os.Exit(1)
	}
	if l := len(hosts); l > 0 {
		fmt.Println(hosts[r.Intn(l)])
	}
}
