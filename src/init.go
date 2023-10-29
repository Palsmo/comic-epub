package src

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	t "github.com/palsmo/comic-epub/src/types"

	"github.com/go-shiori/go-epub"
	color "github.com/mazznoer/csscolorparser"
	"golang.org/x/text/language"
)

// init available options/flags on command
func initOptions(e *epub.Epub, w io.Writer) ([]t.Option, map[string]t.Option) {
	// options
	o1 := t.NewOption("t", 1, []string{"specify the book's title"})
	o2 := t.NewOption("a", 1, []string{"specify the book's author"})
	o3 := t.NewOption("d", 1, []string{"specify the book's description"})
	o4 := t.NewOption("c", -1, []string{"specify the book's chapters, args - P:[TITLE]"})
	o5 := t.NewOption("l", 1, []string{"specify the book's language, ISO 639-[2,3] format"})
	o6 := t.NewOption("bg", 1, []string{"specify the book's page background color,", "needs to be valid css color (e-reader may override)"})
	o7 := t.NewOption("out", 1, []string{"specify the output path, default is current dir"})
	// actions
	o1.Action = func(args ...string) int { e.SetTitle(args[0]); return 0 }
	o2.Action = func(args ...string) int { e.SetAuthor(args[0]); return 0 }
	o3.Action = func(args ...string) int { e.SetDescription(args[0]); return 0 }
	o4.Action = func(args ...string) int {
		f := true
		for _, a := range args {
			s := strings.SplitN(a, ":", 2)
			if len(s) != 2 {
				return errMsg(w, fmt.Sprintf("invalid argument '%s' given to flag -- '%s'", a, o4.Flag))
			}
			n, err := strconv.Atoi(s[0])
			if err != nil || n <= 0 {
				return errMsg(w, fmt.Sprintf("invalid page number '%s' in arg -- '%s'", s[0], a))
			}
			if f {
				pageFirstChapter = n
				f = false
			}
			if v, ok := chapInfo[n]; ok {
				return errMsg(w, fmt.Sprintf("page number '%d' already set to chapter -- '%s'", n, v))
			}
			chapInfo[n] = s[1]
		}
		return 0
	}
	o5.Action = func(args ...string) int {
		if _, err := language.ParseBase(args[0]); err != nil {
			return errMsg(w, fmt.Sprintf("invalid language code -- '%s'", args[0]))
		}
		e.SetLang(args[0])
		return 0
	}
	o6.Action = func(args ...string) int {
		if _, err := color.Parse(args[0]); err != nil {
			return errMsg(w, fmt.Sprintf("invalid background color -- '%s'", args[0]))
		}
		bgColor = args[0]
		return 0
	}
	o7.Action = func(args ...string) int {
		if _, err := os.Stat(args[0]); err != nil {
			return errMsg(w, fmt.Sprintf("invalid output path -- '%s'", args[0]))
		}
		outPath = args[0]
		return 0
	}
	// end
	return []t.Option{o1, o2, o3, o4, o5, o6, o7},
		map[string]t.Option{"t": o1, "a": o2, "d": o3, "c": o4, "l": o5, "bg": o6, "out": o7}
}
