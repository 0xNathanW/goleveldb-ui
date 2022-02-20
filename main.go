package main

import (
	"flag"
	"fmt"
	"log"
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
	uiOpts    ui.UiOpts
)

func main() {

	f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("error opening file:", err)
	}
	log.SetOutput(f)

	flag.Parse()
	fmtOpt, ok := formatMap[formatOpt]
	if !ok {
		fmt.Printf("invalid format: %v\n", formatOpt)
		os.Exit(30)
	}
	uiOpts.Format = fmtOpt

	app := ui.NewUI(dbPath, &uiOpts)
	app.Run()

}

func init() {

	fmtHelp := "How to format the output. Valid options are:\n" +
		"\thex - hexadecimal\n" +
		"\tstring\n" +
		"\tnum - integer\n" +
		"\tbin - binary\n"

	flag.StringVar(&formatOpt, "fmt", "hex", fmtHelp)
	flag.StringVar(&dbPath, "db", "", "Path to database. Required.")
	flag.IntVar(&uiOpts.Max, "max", 5, "Max keys per page.")
}
