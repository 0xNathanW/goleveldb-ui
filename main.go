package main

import (
	"log"
	"os"
)

func main() {

	dbPath := os.Args[1]
	if !verifyPath(dbPath) {
		log.Fatal("Invalid database path")
	}

}

func verifyPath(dbPath string) bool {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false
	}
	return true
}
