package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fordhurley/thorf"
)

func main() {
	var input io.Reader = os.Stdin

	if len(os.Args) == 2 {
		// Read program from file named by first argument
		filename := os.Args[1]
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		input = f
	} else if len(os.Args) == 3 && os.Args[1] == "-e" {
		// Execute program provided as argument:
		input = strings.NewReader(os.Args[2])
	}

	m := thorf.NewMachine(os.Stdout)
	err := m.Eval(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println()
}
