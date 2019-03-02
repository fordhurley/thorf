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

type Stack []int

func (s *Stack) Push(x int) {
	*s = append(*s, x)
}

func (s *Stack) Pop() int {
	x := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return x
}

func (s *Stack) Peek() int {
	return (*s)[len(*s)-1]
}

func (s *Stack) Size() int {
	return len(*s)
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

// Operation modifies the stack or returns an error.
type Operation func(*Stack) error

func add(s *Stack) error {
	if s.Size() < 2 {
		return fmt.Errorf("need two values to add")
	}
	a := s.Pop()
	b := s.Pop()
	s.Push(a + b)
	return nil
}

func subtract(s *Stack) error {
	if s.Size() < 2 {
		return fmt.Errorf("need two values to subtract")
	}
	a := s.Pop()
	b := s.Pop()
	s.Push(b - a)
	return nil
}

func multiply(s *Stack) error {
	if s.Size() < 2 {
		return fmt.Errorf("need two values to multiply")
	}
	a := s.Pop()
	b := s.Pop()
	s.Push(a * b)
	return nil
}

func divide(s *Stack) error {
	if s.Size() < 2 {
		return fmt.Errorf("need two values to divide")
	}
	a := s.Pop()
	if a == 0 {
		return fmt.Errorf("divide by zero")
	}
	b := s.Pop()
	s.Push(b / a)
	return nil
}

func duplicate(s *Stack) error {
	if s.Size() < 1 {
		return fmt.Errorf("need a value to dup")
	}
	s.Push(s.Peek())
	return nil
}

func drop(s *Stack) error {
	if s.Size() < 1 {
		return fmt.Errorf("need a value to drop")
	}
	s.Pop()
	return nil
}

func swap(s *Stack) error {
	if s.Size() < 2 {
		return fmt.Errorf("need two values to swap")
	}
	a := s.Pop()
	b := s.Pop()
	s.Push(a)
	s.Push(b)
	return nil
}

func over(s *Stack) error {
	if s.Size() < 2 {
		return fmt.Errorf("need two values to copy over")
	}
	a := s.Pop()
	b := s.Peek()
	s.Push(a)
	s.Push(b)
	return nil
}
