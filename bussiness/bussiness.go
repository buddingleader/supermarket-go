package bussiness

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"supermarket-go/common"
	db "supermarket-go/local/leveldb"
	"supermarket-go/log"

	"github.com/syndtr/goleveldb/leveldb/util"
)

// 对象键值
const (
	ACCOUNTBOOKNAME = "accountBook" //账簿名称
	GOODSNAME       = "goods"       //商品库名称
)

// accountBook 账目
type accountBook struct {
	BarCode int64   `json:"barCode"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	Time    string  `json:"time"`
}

func (a *accountBook) String() string {
	sb := bytes.Buffer{}
	sb.WriteString(" 条形码：")
	sb.WriteString(common.Int64ToString(a.BarCode))
	sb.WriteString(" 出售商品：")
	sb.WriteString(a.Name)
	sb.WriteString(" 出售价格：")
	sb.WriteString(common.Float64ToString(a.Price))
	sb.WriteString(" 出售时间：")
	sb.WriteString(a.Time)
	return sb.String()
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
func saveAcountBook(date string, barcode int64, price float64) bool {
	allabs := make(map[string][]accountBook)                 // 日期做键值,每天都是一个小账本
	abs := []accountBook{}                                   // 新建当天账本
	if err := db.Get(ACCOUNTBOOKNAME, &allabs); err == nil { //取得账簿
		abs = allabs[date] //已有当天账本则取值当天账本
	}
	goods := make(map[int64]good) //条形码做键值
	good := good{Name: "默认商品"}
	if err := db.Get(GOODSNAME, &goods); err == nil { //取得商品库
		good = goods[barcode] //取得当前卖出的商品
	}
	ab := accountBook{barcode, good.Name, price, common.GetTimeStamp()} //初始化账目
	abs = append(abs, ab)                                               //加入账本
	allabs[date] = abs                                                  //存入账簿
	if db.Put(ACCOUNTBOOKNAME, allabs) {
		log.InfoLog("存入数据库成功！", ab.String())
	}
	return true
}

func showDataBase() {
	log.InfoLog("开始输出数据库...")
	// iter := db.Database.NewIterator(nil, nil)
	iter := db.Database.NewIterator(util.BytesPrefix([]byte(ACCOUNTBOOKNAME)), nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		var abs map[string][]accountBook
		if err := json.Unmarshal(value, &abs); err != nil {
			db.Database.Delete([]byte(ACCOUNTBOOKNAME), nil)
			log.ErrorLog("反格式化数据错误,value：", abs, ".开始执行删除")
		}
		log.InfoLog(string(key), "账簿：")
		var sumPrice float64
		for i, v := range abs {
			sumPrice = v[0].Price
			log.InfoLog("日期：", i, v[0].String())
			for index := 1; index < len(v); index++ {
				log.InfoLog("               ", v[index].String())
				sumPrice += v[index].Price
			}
			log.InfoLog("                                                      总  价 :", common.Float64ToString(sumPrice))

		}
	}
	iter.Release()
	// err = iter.Error()
}
