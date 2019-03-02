package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

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
	}

	m := thorf.NewMachine()

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		err := m.Eval(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	err := scanner.Err()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(m.Stack())
}
