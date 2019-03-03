package thorf

import (
	"fmt"
	"io"
)

// An Operation modifies the stack or returns an error.
type Operation func(*Stack) error

func add(s *Stack) error {
	if s.Size() < 2 {
		return fmt.Errorf("need two values to add")
	}
	a := s.Pop()
	b := s.Pop()
	s.Push(b + a)
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
	s.Push(b * a)
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

func print(w io.Writer) Operation {
	return func(s *Stack) error {
		if s.Size() < 1 {
			return fmt.Errorf("need one value to print")
		}
		x := s.Pop()
		_, err := fmt.Fprint(w, x, " ")
		return err
	}
}

// printStack does NOT consume the stack, unlike most of the other operations.
func printStack(w io.Writer) Operation {
	return func(s *Stack) error {
		for _, x := range *s {
			_, err := fmt.Fprint(w, x, " ")
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func emit(w io.Writer) Operation {
	return func(s *Stack) error {
		if s.Size() < 1 {
			return fmt.Errorf("need one value to emit")
		}
		x := s.Pop()
		_, err := fmt.Fprint(w, string(rune(x)))
		return err
	}
}
