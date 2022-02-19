package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {

	db, err := leveldb.OpenFile("./test-db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 100; i++ {
		putRandomBytes(db, int64(i+2c))
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
