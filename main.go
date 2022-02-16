package main

import (
	"log"
	"os"

	"github.com/0xNathanW/goleveldb-ui/ui"
)

func main() {

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	log.SetOutput(file)

	path := os.Args[1]
	app := ui.NewUI(path)
	app.Run()

}
