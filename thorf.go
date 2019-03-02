package thorf

import (
	"fmt"
	"strconv"
	"strings"
)

// Machine is a Forth interpreter with state.
type Machine struct {
	stack *Stack
	dict  map[string]Operation
}

// NewMachine returns a Machine with an empty stack and the default dictionary.
func NewMachine() *Machine {
	return &Machine{
		stack: &Stack{},
		dict: map[string]Operation{
			"+":    add,
			"-":    subtract,
			"*":    multiply,
			"/":    divide,
			"dup":  duplicate,
			"drop": drop,
			"swap": swap,
			"over": over,
		},
	}
}

// Eval a single line of input.
func (m *Machine) Eval(input string) error {
	words := strings.Fields(input)

	if words[0] == ":" {
		// Add a user defined word:
		word := words[1]

		_, err := strconv.Atoi(word)
		if err == nil {
			return fmt.Errorf("cannot redefine numbers")
		}

		var definition []string
		for _, w := range words[2:] {
			if w == ";" {
				break
			}
			definition = append(definition, w)
		}

		m.dict[strings.ToLower(word)] = defineOperation(m.dict, definition)

		return nil
	}

	return eval(m.stack, m.dict, words)
}

// Stack returns the current state of the Machine stack.
func (m *Machine) Stack() []int {
	return *m.stack
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
