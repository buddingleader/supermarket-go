package leveldb

import (
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

// Database
var (
	Database *leveldb.DB
)

func init() {
	db, err := leveldb.OpenFile("/db", nil)
	if err != nil {
		panic("打开数据库错误...")
	}
	Database = db
	fmt.Println("数据库初始化成功！")
	// defer Database.Close()
}

// Put 存储key-value
func Put(key string, value interface{}) bool {
	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	err = Database.Put([]byte(key), v, nil)
	if err != nil {
		panic(err)
	}
	return true
}

// Get 根据key和value类型获取value
func Get(key string, value interface{}) error {
	v, err := Database.Get([]byte(key), nil)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(v, value); err != nil {
		panic(err)
	}
	return nil
}
