package bussiness

import (
	"bufio"
	"encoding/json"
	"os"
	"supermarket-go/common"
	db "supermarket-go/local/leveldb"
	"supermarket-go/log"

	"github.com/syndtr/goleveldb/leveldb/util"
)

type accountBook struct {
	BarCode int64   `json:"barCode"`
	Price   float32 `json:"price"`
	Time    string  `json:"time"`
}

// OpenBussiness 开始营业
func OpenBussiness() {
	log.InfoLog("开始营业...")
	key := common.GetDate()
	for {
		input := readConsole()
		log.InfoLog("input:", input)
		switch input {
		case "s":
			showDataBase()
		default:
			account(key, input)
		}
	}
}

// readConsole 读取控制台输入,并去掉输入的最后两个字节
func readConsole() string {
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		log.ErrorLog("读取控制台输入错误:", err)
	}
	return input[:len(input)-2]
}

// account 记账
func account(key string, input string) {
	if barcode, isBarcode := checkBarcode(input); isBarcode {
		log.InfoLog("标准条形码输入，请输入价格：")
		for {
			input1 := readConsole()
			if price, isPrice := checkPrice(input1); isPrice {
				saveAcountBook(key, barcode, price)
				break
			} else {
				log.InfoLog("请输入正确的价格：")
			}
		}
	} else {
		log.InfoLog("非条形码输入，尝试匹配价格...")
		if price, isPrice := checkPrice(input); isPrice {
			saveAcountBook(key, common.ToInt64(key), price)
		} else {
			log.InfoLog("请输入正确的价格：")
		}
	}
}

// saveAcountBook 记录账本
func saveAcountBook(date string, input int64, price float32) bool {
	allabs := make(map[string][]accountBook)              // 日期做键值,每天都是一个小账本
	abs := []accountBook{}                                // 新建当天账本
	if err := db.Get("acountBook", &allabs); err == nil { //取得账簿
		abs = allabs[date] //已有当天账本则取值当天账本
	}
	ab := accountBook{input, price, common.GetTimeStamp()} //初始化账目
	abs = append(abs, ab)                                  //加入账本
	allabs[date] = abs                                     //存入账簿
	if db.Put("acountBook", allabs) {
		log.InfoLog("存入数据库成功！", ab)
	}
	return true
}

func showDataBase() {
	log.InfoLog("开始输出数据库...")
	// iter := db.Database.NewIterator(nil, nil)
	iter := db.Database.NewIterator(util.BytesPrefix([]byte("acountBook")), nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		var abs map[string][]accountBook
		if err := json.Unmarshal(value, &abs); err != nil {
			db.Database.Delete([]byte("acountBook"), nil)
			log.ErrorLog("反格式化数据错误,value：", abs, ".开始执行删除")
		}
		log.InfoLog(string(key), abs)
	}
	iter.Release()
	// err = iter.Error()
}
