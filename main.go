package main

import (
	"os"

	"github.com/0xNathanW/goleveldb-ui/ui"
)

func main() {

	path := os.Args[1]
	app := ui.NewUI(path)
	app.Run()

}
