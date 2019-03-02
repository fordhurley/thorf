package thorf

import (
	"fmt"
	"strconv"
	"strings"
)

// Eval evaluates the input statements and returns the end state of the stack.
func Eval(input []string) ([]int, error) {
	var stack Stack
	dict := map[string]Operation{
		"+":    add,
		"-":    subtract,
		"*":    multiply,
		"/":    divide,
		"dup":  duplicate,
		"drop": drop,
		"swap": swap,
		"over": over,
	}

	for _, line := range input {
		words := strings.Fields(line)

		if words[0] == ":" {
			// Add a user defined word:
			word := words[1]

			_, err := strconv.Atoi(word)
			if err == nil {
				return nil, fmt.Errorf("cannot redefine numbers")
			}

			var definition []string
			for _, w := range words[2:] {
				if w == ";" {
					break
				}
				definition = append(definition, w)
			}

			dict[strings.ToLower(word)] = defineOperation(dict, definition)

			continue
		}

		err := eval(&stack, dict, words)
		if err != nil {
			return nil, err
		}
	}

	return stack, nil
}

func defineOperation(dict map[string]Operation, words []string) Operation {
	// A user defined word should operate on a snapshot of the dictionary
	// at the time it was defined.
	clone := make(map[string]Operation, len(dict))
	for w, op := range dict {
		clone[w] = op
	}
	return func(s *Stack) error {
		return eval(s, clone, words)
	}
}

func eval(stack *Stack, dict map[string]Operation, words []string) error {
	for _, word := range words {
		n, err := strconv.Atoi(word)
		if err == nil {
			stack.Push(n)
			continue
		}

		op, ok := dict[strings.ToLower(word)]
		if !ok {
			return fmt.Errorf("unknown word: %q", word)
		}

		err = op(stack)
		if err != nil {
			return err
		}
	}
	return nil
}
