package thorf

import "fmt"

func ExampleStack() {
	var s Stack
	fmt.Println("size:", s.Size())

	s.Push(42)
	fmt.Println("size:", s.Size())

	s.Push(13)
	fmt.Println("size:", s.Size())

	fmt.Println("peek:", s.Peek())
	fmt.Println("size:", s.Size())

	fmt.Println("pop:", s.Pop())
	fmt.Println("size:", s.Size())

	fmt.Println("pop:", s.Pop())
	fmt.Println("size:", s.Size())

	s.Push(108)
	fmt.Println("size:", s.Size())
	fmt.Println("pop:", s.Pop())
	fmt.Println("size:", s.Size())

	// Output:
	// size: 0
	// size: 1
	// size: 2
	// peek: 13
	// size: 2
	// pop: 13
	// size: 1
	// pop: 42
	// size: 0
	// size: 1
	// pop: 108
	// size: 0
}
