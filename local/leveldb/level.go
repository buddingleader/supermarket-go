package leveldb

import (
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/wangff15386/supermarket-go/conf"
)

// Database
var (
	Database *leveldb.DB
)

// Initial the leveldb
func Initial() {
	db, err := leveldb.OpenFile(conf.Config.LevelPath, nil)
	if err != nil {
		panic("打开数据库错误...")
	}
	Database = db
	fmt.Println("数据库初始化成功！")
	// defer Database.Close()
}

// Put 存储key-value
func Put(key string, value interface{}) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = Database.Put([]byte(key), v, nil); err != nil {
		return err
	}
	return nil
}

// Get 根据key和value类型获取value
func Get(key string, value interface{}) error {
	v, err := Database.Get([]byte(key), nil)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(v, value); err != nil {
		return err
	}
	return nil
}
