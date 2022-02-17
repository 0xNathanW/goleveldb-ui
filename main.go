package main

import (
	"fmt"
	"os"

	"github.com/0xNathanW/goleveldb-ui/ui"
	flag "github.com/alecthomas/kong"
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
	//fmt.Println("path: ", path)

	// app := ui.NewUI(path, &ui.UiOpts{
	// 	Format: formatMap[formatOpt],
	// })
	// app.Run()

}

func init() {
	flag.StringVar(&formatOpt, "fmt", "hex", "format to display keys and values")
	flag.StringVar(&dbPath, "db", "", "path to database")
}
