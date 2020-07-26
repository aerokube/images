package main

import (
	"github.com/aerokube/selenoid-images/cmd"
)

//go:generate pkger -include /static -o build

func main() {
	cmd.Execute()
}

