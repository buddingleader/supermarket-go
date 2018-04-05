package leveldb

import (
	"supermarket-go/log"

	"github.com/syndtr/goleveldb/leveldb"
)

// Datebase 数据库引用
var Datebase *leveldb.DB

func init() {
	db, err := leveldb.OpenFile("./local/leveldb/db", nil)
	if err != nil {
		log.ErrorLog("打开数据库错误...")
	}
	Datebase = db
	log.InfoLog("数据库初始化成功！")
	// defer Datebase.Close()
}
