package thorf

import (
	"fmt"
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
	lexer := NewLexer(strings.NewReader(input))

	for lexer.Scan() {
		token := lexer.Token()

		if token.Type == Def {
			// Add a user defined word.

			if !lexer.Scan() {
				err := lexer.Err()
				if err != nil {
					return err
				}
				return fmt.Errorf("invalid word definition")
			}

			t := lexer.Token()
			if t.Type != Word {
				return fmt.Errorf("invalid word definition")
			}
			name := t.word

			var definition []Token
			for lexer.Scan() {
				t = lexer.Token()
				if t.Type == End {
					break
				}
				definition = append(definition, t)
			}

			m.dict[name] = defineOperation(m.dict, definition)

			continue
		}

		err := eval(m.stack, m.dict, token)
		if err != nil {
			return err
		}
	}

	return lexer.Err()
}

// Stack returns the current state of the Machine stack.
func (m *Machine) Stack() []int {
	return *m.stack
}

func defineOperation(dict map[string]Operation, definition []Token) Operation {
	// A user defined word should operate on a snapshot of the dictionary
	// at the time it was defined.
	clone := make(map[string]Operation, len(dict))
	for w, op := range dict {
		clone[w] = op
	}
	return func(s *Stack) error {
		return eval(s, clone, definition...)
	}
}

func eval(stack *Stack, dict map[string]Operation, tokens ...Token) error {
	for _, token := range tokens {
		switch token.Type {
		case Word:
			op, ok := dict[token.word]
			if !ok {
				return fmt.Errorf("unknown word: %q", token.word)
			}
			err := op(stack)
			if err != nil {
				return err
			}
		case Num:
			stack.Push(token.num)
		default:
			panic("unexpected token: " + token.String())
		}
	}
	return nil
}
