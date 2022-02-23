package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {

	f, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("error opening file:", err)
	}
	log.SetOutput(f)

	db, err := leveldb.OpenFile("./test-db", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < 1000; i++ {
		_ = db.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)), nil)
	}

	//iterateDB(db)
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
	k := make([]byte, 20)
	v := make([]byte, 100)
	_, err := rand.Read(k)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(seed * 2)
	_, err = rand.Read(v)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(k, v)
	err = db.Put(k, v, nil)
	if err != nil {
		log.Fatal(err)
	}
}
