package thorf

import (
	"bytes"
	"strings"
	"testing"
)

// These tests were borrowed from the excellent exercism.io "forth" exercise:
// https://github.com/exercism/go/tree/5446524b6/exercises/forth

func runTest(input string) (string, error) {
	var buf bytes.Buffer
	m := NewMachine(&buf)

	err := m.Eval(strings.NewReader(input))
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func TestMachine(t *testing.T) {
	for _, tg := range testGroups {
		for _, tc := range tg.tests {
			t.Run(tg.group+"--"+tc.description, func(t *testing.T) {
				output, err := runTest(tc.input)
				if err != nil {
					if !tc.err {
						t.Fatalf("runTest(%#v) expected %q, got an error: %q", tc.input, tc.expected, err)
					}
					return
				}

				if tc.err {
					t.Fatalf("runTest(%#v) expected an error, got %q", tc.input, output)
				}

				if output != tc.expected {
					t.Fatalf("runTest(%#v) expected %q, got %q", tc.input, tc.expected, output)
				}
			})
		}
	}
}

func BenchmarkMachine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tg := range testGroups {
			for _, tc := range tg.tests {
				runTest(tc.input)
			}
		}
	}
}

type testGroup struct {
	group string
	tests []testCase
}

type testCase struct {
	description string
	input       string
	expected    string
	err         bool
}

var testGroups = []testGroup{
	{
		group: "parsing and numbers",
		tests: []testCase{
			{
				"numbers just get pushed onto the stack",
				"1 2 3 4 5 .s",
				"1 2 3 4 5 ",
				false,
			},
		},
	},
	{
		group: "addition",
		tests: []testCase{
			{
				"can add two numbers",
				"1 2 + .s",
				"3 ",
				false,
			},
			{
				"errors if there is nothing on the stack",
				"+",
				"",
				true,
			},
			{
				"errors if there is only one value on the stack",
				"1 +",
				"",
				true,
			},
		},
	},
	{
		group: "subtraction",
		tests: []testCase{
			{
				"can subtract two numbers",
				"3 4 - .s",
				"-1 ",
				false,
			},
			{
				"errors if there is nothing on the stack",
				"-",
				"",
				true,
			},
			{
				"errors if there is only one value on the stack",
				"1 -",
				"",
				true,
			},
		},
	},
	{
		group: "multiplication",
		tests: []testCase{
			{
				"can multiply two numbers",
				"2 4 * .s",
				"8 ",
				false,
			},
			{
				"errors if there is nothing on the stack",
				"*",
				"",
				true,
			},
			{
				"errors if there is only one value on the stack",
				"1 *",
				"",
				true,
			},
		},
	},
	{
		group: "division",
		tests: []testCase{
			{
				"can divide two numbers",
				"12 3 / .s",
				"4 ",
				false,
			},
			{
				"performs integer division",
				"8 3 / .s",
				"2 ",
				false,
			},
			{
				"errors if dividing by zero",
				"4 0 /",
				"",
				true,
			},
			{
				"errors if there is nothing on the stack",
				"/",
				"",
				true,
			},
			{
				"errors if there is only one value on the stack",
				"1 /",
				"",
				true,
			},
		},
	},
	{
		group: "combined arithmetic",
		tests: []testCase{
			{
				"addition and subtraction",
				"1 2 + 4 - .s",
				"-1 ",
				false,
			},
			{
				"multiplication and division",
				"2 4 * 3 / .s",
				"2 ",
				false,
			},
		},
	},
	{
		group: "dup",
		tests: []testCase{
			{
				"copies a value on the stack",
				"1 dup .s",
				"1 1 ",
				false,
			},
			{
				"copies the top value on the stack",
				"1 2 dup .s",
				"1 2 2 ",
				false,
			},
			{
				"errors if there is nothing on the stack",
				"dup",
				"",
				true,
			},
		},
	},
	{
		group: "drop",
		tests: []testCase{
			{
				"removes the top value on the stack if it is the only one",
				"1 drop .s",
				"",
				false,
			},
			{
				"removes the top value on the stack if it is not the only one",
				"1 2 drop .s",
				"1 ",
				false,
			},
			{
				"errors if there is nothing on the stack",
				"drop",
				"",
				true,
			},
		},
	},
	{
		group: "swap",
		tests: []testCase{
			{
				"swaps the top two values on the stack if they are the only ones",
				"1 2 swap .s",
				"2 1 ",
				false,
			},
			{
				"swaps the top two values on the stack if they are not the only ones",
				"1 2 3 swap .s",
				"1 3 2 ",
				false,
			},
			{
				"errors if there is nothing on the stack",
				"swap",
				"",
				true,
			},
			{
				"errors if there is only one value on the stack",
				"1 swap",
				"",
				true,
			},
		},
	},
	{
		group: "over",
		tests: []testCase{
			{
				"copies the second element if there are only two",
				"1 2 over .s",
				"1 2 1 ",
				false,
			},
			{
				"copies the second element if there are more than two",
				"1 2 3 over .s",
				"1 2 3 2 ",
				false,
			},
			{
				"errors if there is nothing on the stack",
				"over",
				"",
				true,
			},
			{
				"errors if there is only one value on the stack",
				"1 over",
				"",
				true,
			},
		},
	},
	{
		group: "user-defined words",
		tests: []testCase{
			{
				"can consist of built-in words",
				": dup-twice dup dup ; 1 dup-twice .s",
				"1 1 1 ",
				false,
			},
			{
				"execute in the right order",
				": countup 1 2 3 ; countup .s",
				"1 2 3 ",
				false,
			},
			{
				"can override other user-defined words",
				": foo dup ; : foo dup dup ; 1 foo .s",
				"1 1 1 ",
				false,
			},
			{
				"can override built-in words",
				": swap dup ; 1 swap .s",
				"1 1 ",
				false,
			},
			{
				"can override built-in operators",
				": + * ; 3 4 + .s",
				"12 ",
				false,
			},
			{
				"can use different words with the same name",
				": foo 5 ; : bar foo ; : foo 6 ; bar foo .s",
				"5 6 ",
				false,
			},
			{
				"can define word that uses word with the same name",
				": foo 10 ; : foo foo 1 + ; foo .s",
				"11 ",
				false,
			},
			{
				"cannot redefine numbers",
				": 1 2 ;",
				"",
				true,
			},
			{
				"errors if executing a non-existent word",
				"foo",
				"",
				true,
			},
		},
	},
	{
		group: "case-insensitivity",
		tests: []testCase{
			{
				"DUP is case-insensitive",
				"1 DUP Dup dup .s",
				"1 1 1 1 ",
				false,
			},
			{
				"DROP is case-insensitive",
				"1 2 3 4 DROP Drop drop .s",
				"1 ",
				false,
			},
			{
				"SWAP is case-insensitive",
				"1 2 SWAP 3 Swap 4 swap .s",
				"2 3 4 1 ",
				false,
			},
			{
				"OVER is case-insensitive",
				"1 2 OVER Over over .s",
				"1 2 1 2 1 ",
				false,
			},
			{
				"user-defined words are case-insensitive",
				": foo dup ; 1 FOO Foo foo .s",
				"1 1 1 1 ",
				false,
			},
			{
				"definitions are case-insensitive",
				": SWAP DUP Dup dup ; 1 swap .s",
				"1 1 1 1 ",
				false,
			},
		},
	},
	{
		group: "output",
		tests: []testCase{
			{
				". prints the value as a number",
				"1 .",
				"1 ",
				false,
			},
			{
				". prints only the last value as a number",
				"1 2 3 .",
				"3 ",
				false,
			},
			{
				". consumes the value",
				"1 2 3 . .s",
				"3 1 2 ",
				false,
			},
			{
				".s prints the whole stack",
				"1 2 3 .s",
				"1 2 3 ",
				false,
			},
			{
				".s does not consume the stack",
				"1 .s .",
				"1 1 ",
				false,
			},
			{
				"EMIT prints one ASCII character",
				"42 EMIT",
				"*",
				false,
			},
			{
				"EMIT consumes the value",
				"42 EMIT .s",
				"*",
				false,
			},
			{
				". errors if there is nothing on the stack",
				".",
				"",
				true,
			},
			{
				".s can print an empty stack",
				".s",
				"",
				false,
			},
			{
				"EMIT errors if there is nothing on the stack",
				"EMIT",
				"",
				true,
			},
		},
	},
}
