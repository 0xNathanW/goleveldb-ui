package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/0xNathanW/goleveldb-ui/ui"
)

var (
	formatMap = map[string]uint{
		"hex":    ui.Hex,
		"string": ui.Str,
		"num":    ui.Num,
		"bin":    ui.Bin,
	}
)

var (
	dbPath    string
	formatOpt string
)

func main() {
	flag.Parse()
	fmtOpt, ok := formatMap[formatOpt]
	if !ok {
		fmt.Printf("invalid format: %v\n", formatOpt)
		os.Exit(1)
	}

	fmt.Println("format: ", fmtOpt)
	fmt.Println("path: ", dbPath)

	app := ui.NewUI(dbPath, &ui.UiOpts{
		Format: formatMap[formatOpt],
	})
	app.Run()

}

func init() {

	fmtHelp := "How to format the output. Valid options are:\n" +
		"\thex - hexadecimal\n" +
		"\tstring\n" +
		"\tnum - integer\n" +
		"\tbin - binary\n"

	flag.StringVar(&formatOpt, "fmt", "hex", fmtHelp)
	flag.StringVar(&dbPath, "db", "", "Path to database.\n")
}
