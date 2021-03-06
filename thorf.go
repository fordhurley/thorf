package thorf

import (
	"fmt"
	"io"
)

// Machine is a Forth interpreter with state.
type Machine struct {
	stack *Stack
	dict  map[string]Operation
}

// NewMachine returns a Machine with an empty stack and the default dictionary.
func NewMachine(w io.Writer) *Machine {
	return &Machine{
		stack: &Stack{},
		dict: map[string]Operation{
			// Arithmetic:
			"+": add,
			"-": subtract,
			"*": multiply,
			"/": divide,
			// Stack manipulations:
			"dup":  duplicate,
			"drop": drop,
			"swap": swap,
			"over": over,
			// Output:
			".":    print(w),
			".s":   printStack(w),
			"emit": emit(w),
		},
	}
}

// Eval evaluates instructions read from r.
func (m *Machine) Eval(r io.Reader) error {
	lexer := NewLexer(r)

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
