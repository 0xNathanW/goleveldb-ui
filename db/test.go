package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {

	db, err := leveldb.OpenFile("./test-db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 100; i++ {
		if err := db.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)), nil); err != nil {
			log.Fatal(err)
		}
	}

	iterateDB(db)
}

func iterateDB(db *leveldb.DB) {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		fmt.Printf("Key: %x\tValue: %x\n", iter.Key(), iter.Value())
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Fatal(err)
	}
}

func flipFlop(db *leveldb.DB) {
	iter := db.NewIterator(nil, nil)
	for i := 0; i < 10; i++ {
		fmt.Println("Next")
		iter.Next()
		fmt.Printf("Key: %x\tValue: %x\n", iter.Key(), iter.Value())
		fmt.Println("Next")
		iter.Next()
		fmt.Printf("Key: %x\tValue: %x\n", iter.Key(), iter.Value())
		fmt.Println("Prev")
		iter.Prev()
		fmt.Printf("Key: %x\tValue: %x\n", iter.Key(), iter.Value())
	}
}

// func reverseIterateDB(db *leveldb.DB) {
// 	iter := db.NewIterator(nil, nil)

// 	for iter.Last() {

// }

func putRandomBytes(db *leveldb.DB, seed int64) {
	rand.Seed(seed)
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Put(b, b, nil)
	if err != nil {
		log.Fatal(err)
	}
}
