package src

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	f "path/filepath"

	"github.com/go-shiori/go-epub"
)

func Setup(w io.Writer) int {
	// ready epub struct
	var err error
	e, err = epub.NewEpub("comic_new")
	if err != nil {
		return errMsg(os.Stderr, err.Error())
	}
	// init
	optList, optDict = initOptions(e, w)
	defaultVars()
	// end
	return 0
}

func Run(w io.Writer) int {
	// used to silence output during test
	var ew io.Writer
	if w == io.Discard {
		ew = io.Discard
	} else {
		ew = os.Stderr
	}
	// parse input arguments
	args := os.Args[1:]
	if len(args) == 0 {
		printInfo(w, Cmd, optList)
		return 0
	}
	// css
	data := "data:text/css," +
		"body%7Bbackground-color:" + url.QueryEscape(bgColor) + ";margin:0;text-align:center;%7D" +
		"img%7Bmax-height:100%25;max-width:100%25;%7D"
	css, err := e.AddCSS(data, "style.css")
	if err != nil {
		return errMsg(ew, err.Error())
	}
	// flags
	hasCov := false  // is cover set (first image)
	hasPath := false // is image path argument given
	hasImg := false  // was at least one image retrieved
	// cache
	totN, relN := 0, 0
	// for each given argument
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg[0] == '-' {
			// arg' - handled as option
			o, ok := optDict[arg[1:]]
			if !ok {
				return errMsg(ew, fmt.Sprintf("invalid option -- '%s'", arg[1:]))
			}
			// get arguments for option 'o'
			var o_args []string
			if o.Nargs == -1 {
				for i < len(args)-1 {
					i++ // point to proceeding arg
					if args[i][0] == ',' {
						break
					}
					o_args = append(o_args, args[i])
				}
			} else {
				for j := 0; j < o.Nargs; j++ {
					i++ // point to proceeding arg
					if i >= len(args) {
						break
					}
					o_args = append(o_args, args[i])
				}
			}
			if o.Nargs > len(o_args) || len(o_args) == 0 {
				return errMsg(ew, fmt.Sprintf("too few arguments for option -- '%s'", o.Flag))
			}
			if o.Action(o_args...) == 1 {
				return 1
			}
			continue
		}
		// 'arg' - handled as path
		pth, err := os.Stat(arg)
		if err != nil {
			return errMsg(ew, fmt.Sprintf("invalid path -- '%s'", arg))
		}
		// flag
		hasPath = true
		// when called, put image 'file' into 'e' (epub)
		var makePage = func(file string) int {
			// when 'pth' is directory, 'file' only has basename, has to extend
			p := arg
			if pth.IsDir() {
				p = f.Join(p, file)
			}
			// set name
			var name string
			if len(chapInfo) == 0 {
				name = fmt.Sprintf("%04d", totN)
			} else {
				name = fmt.Sprintf("%04d-%04d", totN, relN)
			}
			name += f.Ext(file)
			// include 'file' in 'e' (epub)
			img, err := e.AddImage(p, name)
			if err != nil {
				return errMsg(ew, err.Error())
			}
			// if cover is not set
			if !hasCov {
				e.SetCover(img, css)
				hasCov = true
				totN++
				return 0
			}
			// chapter title
			title := ""
			if len(chapInfo) > 0 {
				if v, ok := chapInfo[totN+1]; ok {
					title = v
				}
				if _, ok := chapInfo[totN+2]; ok {
					relN = 1
				} else if totN+1 >= pageFirstChapter {
					relN++
				}
			}
			totN++
			// try add section
			body := "<img src=\"" + img + "\" alt=\"unsupported image format\"/>"
			if _, err = e.AddSection(body, title, name+".xhtml", css); err != nil {
				return errMsg(ew, err.Error())
			}
			return 0
		}
		// if argument is path to a directory
		// when called, try put image 'file' into 'e' (epub)
		var h = func(_pth os.DirEntry) int {
			ext := f.Ext(_pth.Name())
			if imgFormats.Contains(ext) {
				hasImg = true
				return makePage(_pth.Name())
			}
			return 0
		}
		if pth.IsDir() {
			sub_pths, err := os.ReadDir(arg)
			if err != nil {
				return errMsg(ew, err.Error())
			}
			// for each sub path in argument specified directory
			for _, sub := range sub_pths {
				if h(sub) == 1 {
					return 1
				}
			}
			continue
		}
		// argument is path to file
		if h(fs.FileInfoToDirEntry(pth)) == 1 {
			return 1
		}
	}
	// check flags
	if !hasPath {
		return errMsg(ew, "missing path to image file(s)")
	}
	if !hasImg {
		fmt.Fprintf(w, "%s: no image(s) found in specified path(s)\n", Cmd)
		return 0
	}
	// write epub file
	if err := e.Write(f.Join(outPath, e.Title()+".epub")); err != nil {
		return errMsg(ew, err.Error())
	}
	// end
	fmt.Fprintf(w, "%s: .epub-file created\n", Cmd)
	return 0
}
