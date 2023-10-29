package src

import (
	"fmt"
	"io"
	"strings"

	t "github.com/palsmo/comic-epub/src/types"
)

// clear all content of map
func clearMap[K comparable, V any](m map[K]V) {
	for k := range m {
		delete(m, k)
	}
}

// handler for error exit
func errMsg(w io.Writer, msg string) int {
	fmt.Fprintf(w, "%s: %s\n", Cmd, msg)
	return 1
}

// max for int
func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// print info about command
func printInfo(w io.Writer, cmd string, opts []t.Option) {
	intro := `
usage :
        comic-epub [OPTIONS]... [DIR] 'or [FILE FILE...]
info  :
        create .epub-file from images
ex    :
        comic-epub . # works

        \ - line continuation
        , - stop read args (when flag has no arg limit)
        default \ same effect without
		
        (standing in dir with images)

        comic-epub \
        -t "The Epic Comic" -a "Stan Lee" \
        -d "Story of epic proportions!" \
        -c 7:"I: The Intro" 42:"II: The Epic" , \ # don't forget ','
        -l "eng" \    # default
        -bg "#000" \  # default
        -out . \      # default
        .
`
	fmt.Fprint(w, intro)
	fmt.Fprintln(w, "flags :   ")
	// for each option
	for _, o := range opts {
		// print first row
		f := "       -" + o.Flag
		fmt.Fprintln(w, f+strings.Repeat(" ", 7-len(o.Flag))+o.Info[0])
		// print other rows
		for _, in := range o.Info[1:] {
			ofs := 10 - len(o.Flag)
			fmt.Fprintln(w, strings.Repeat(" ", len(f))+strings.Repeat(" ", max(0, ofs))+in)
		}
	}
	fmt.Fprintln(w)
}
