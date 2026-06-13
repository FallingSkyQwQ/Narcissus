package main

import (
	"os"

	"github.com/FallingSkyQwQ/Narcissus/cmd/narc/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
