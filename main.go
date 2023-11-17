package main

import (
	"os"

	"github.com/mstorm/sprite/internal/service"
)

func main() {
	// TODO: arguments
	name := "sprites"
	files := os.Args[1:]

	if err := service.Gen(name, files); err != nil {
		panic(err)
	}
}
