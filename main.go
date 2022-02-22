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
	dbPath string
	keyFmt string
	valFmt string
	uiOpts ui.UiOpts
)

func main() {

	f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("error opening file:", err)
	}
	log.SetOutput(f)

	flag.Parse()
	key, ok := formatMap[keyFmt]
	if !ok {
		fmt.Printf("invalid format: %v\n", keyFmt)
		os.Exit(30)
	}
	if key != ui.Str && key != ui.Hex {
		fmt.Println("key format must be string or hex")
		os.Exit(30)
	}

	val, ok := formatMap[valFmt]
	if !ok {
		fmt.Printf("invalid format: %v\n", valFmt)
		os.Exit(30)
	}
	uiOpts.KeyFmt, uiOpts.ValFmt = key, val

	app := ui.NewUI(dbPath, &uiOpts)
	app.Run()

}

func init() {

	flag.StringVar(&keyFmt, "key", "string", "Key format")
	flag.StringVar(&valFmt, "val", "string", "Value format")
	flag.StringVar(&dbPath, "db", "", "Path to database. Required.")
	flag.IntVar(&uiOpts.Max, "max", 5, "Max keys per page.")

}
