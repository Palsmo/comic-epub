package main

import (
	"os"

	s "github.com/palsmo/comic-epub/src"
)

func main() {
	w := os.Stdout
	if s.Setup(w) == 1 {
		os.Exit(1)
	}
	os.Exit(s.Run(w))
}
