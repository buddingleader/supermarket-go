package leveldb

import (
	"encoding/json"
	"supermarket-go/log"

	"github.com/syndtr/goleveldb/leveldb"
)

// Database 数据库引用
var Database *leveldb.DB

func init() {
	db, err := leveldb.OpenFile("/db", nil)
	if err != nil {
		log.ErrorLog("打开数据库错误...")
	}
	Database = db
	log.InfoLog("数据库初始化成功！")
	// defer Database.Close()
}

// Put 存储key-value
func Put(key string, value interface{}) bool {
	v, err := json.Marshal(value)
	if err != nil {
		log.ErrorLog("序列化失败！", err)
	}
	err = Database.Put([]byte(key), v, nil)
	if err != nil {
		log.ErrorLog("存入数据库失败！", err)
		return false
	}
	return true
}

// Get 根据key和value类型获取value
func Get(key string, value interface{}) error {
	v, err := Database.Get([]byte(key), nil)
	if err != nil {
		log.ErrorLog("获取key错误！key:", key, "err:", err)
		return err
	}
	if err := json.Unmarshal(v, value); err != nil {
		log.ErrorLog("反序列化失败，v:", v, "value:", value)
		return err
	}
	return nil
}
