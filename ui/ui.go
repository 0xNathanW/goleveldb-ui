package ui

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type ui struct {
	db *leveldb.DB
}

func NewUI(dbPath string) {

}
