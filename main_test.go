package main

import (
	"io"
	"os"
	"strings"
	"testing"

	s "github.com/palsmo/comic-epub/src"
)

func initArgs(arg string) {
	os.Args = []string{s.Cmd}
	if len(arg) > 0 {
		os.Args = append(os.Args, strings.FieldsFunc(arg, func(r rune) bool {
			return r == ' ' || r == '\n'
		})...)
	}
}

func TestMainArgs(t *testing.T) {
	w := io.Discard // shorthand
	t.Parallel()
	// test cases
	tcs := []struct {
		arg  string
		code int
	}{
		// basic
		{"", 0},   // no args -> valid
		{".", 0},  // current path -> valid
		{"~", 1},  // invalid path -> error
		{"-~", 1}, // unknown flag -> error
		// single argument option
		{"-t", 1},     // too few args -> error
		{"-t ~", 1},   // missing path -> error
		{"-t ~ .", 0}, // has arg and path -> valid
		// // language format
		{"-l ~ .", 1},   // bad parse -> error
		{"-l eng .", 0}, // ok parse -> valid
		// background color format
		{"-bg ~ .", 1},    // bad parse -> error
		{"-bg #000 .", 0}, // ok parse -> valid
		// output path
		{"-out ~ .", 1}, // invalid path -> error
		{"-out . .", 0}, // ok path -> valid
		// multi argument option & chapter format
		{"-c", 1},              // too few args -> error
		{"-c ~", 1},            // invalid format on arg -> error
		{"-c N:~", 1},          // invalid page number -> error
		{"-c 7:~", 1},          // missing path -> error
		{"-c 7:~ 7:~", 1},      // page number already set -> error
		{"-c 7:~ 42:~", 1},     // no stop, missing path -> error
		{"-c 7:~ 42:~ ,", 1},   // with stop, missing path -> error
		{"-c 7:~ 42:~ , .", 0}, // valid arg(s), stop and path -> valid
		// full command
		{"-t ~ -a ~ -d ~ -c 7:~ 42:~ , -l eng -bg #000 -out . .", 0},
	}
	// run each test case
	for _, tc := range tcs {
		t.Run(tc.arg, func(t *testing.T) {
			// init
			initArgs(tc.arg)
			s.Setup(w)
			// try run
			if tc.code != s.Run(w) {
				t.Fail()
			}
		})
	}
}
