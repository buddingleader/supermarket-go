package app

import (
	"fmt"
	"testing"

	"github.com/wangff15386/supermarket-go/conf"
	"github.com/wangff15386/supermarket-go/local/leveldb"
)

func Test_Start(t *testing.T) {
	conf.Initial("../conf/app.conf")
	leveldb.Initial()

	db := leveldb.Database
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key), string(value))
	}
	// Start()
}
