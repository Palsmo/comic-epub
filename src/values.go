package src

import (
	t "github.com/palsmo/comic-epub/src/types"

	"github.com/go-shiori/go-epub"
)

const Cmd = "comic-epub" // script name

var (
	e          *epub.Epub          // epub struct
	optDict    map[string]t.Option // option map
	optList    []t.Option          // option list
	imgFormats = t.NewSet[string]( // valid image types
		".apng",
		".avif",
		".gif",
		".jpeg",
		".jpg",
		".png",
		".svg",
		".webp",
	)
)
var (
	chapInfo         map[int]string // map of page-title pairs
	lang             string         // 2- or 3-letter ISO 639 language code
	bgColor          string         // page background color
	outPath          string         // output path for epub file
	pageFirstChapter int            // page for first chapter (relative chapter is part of file's name inside epub)
)

// set values to default
func defaultVars() {
	chapInfo = make(map[int]string)
	lang = "eng"
	bgColor = "#000"
	outPath = "."
	pageFirstChapter = 0
}
