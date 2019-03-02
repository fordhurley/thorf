package thorf

// Stack is a very simple implementation of a stack data structure for ints.
type Stack []int

// Push a value onto the Stack.
func (s *Stack) Push(x int) {
	*s = append(*s, x)
}

// Pop removes and returns the last value that was pushed to the Stack. Panics
// if called on an empty Stack.
func (s *Stack) Pop() int {
	x := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return x
}

// Peek returns the last value from the Stack, without removing it.
func (s *Stack) Peek() int {
	return (*s)[len(*s)-1]
}

// Size returns the number of items in the Stack.
func (s *Stack) Size() int {
	return len(*s)
}
