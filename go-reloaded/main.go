package main

import (
	"errors"
	"fmt"
	"os"
)

func run(args []string) error {
	if len(args) != 2 {
		return errors.New("usage: go run . <input.txt> <output.txt>")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	result := processText(string(data))
	return os.WriteFile(args[1], []byte(result), 0o644)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
