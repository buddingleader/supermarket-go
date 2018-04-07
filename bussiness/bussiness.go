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

// 账目输出
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

// ToString 账簿输出
func ToString(allabs map[string][]accountBook) string {
	sb := bytes.Buffer{}
	sb.WriteString(string(ACCOUNTBOOKNAME))
	sb.WriteString(" 账簿：]\n")
	var sumPrice float64
	for i, abs := range allabs {
		if len(abs) > 0 {
			sumPrice = abs[0].Price
			sb.WriteString("<--日期：")
			sb.WriteString(i)
			sb.WriteString("-->\n")
			sb.WriteString("     —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— \n")
			sb.WriteString("    |")
			sb.WriteString(abs[0].String())
			sb.WriteString("|\n")
			for index := 1; index < len(abs); index++ {
				sb.WriteString("    |")
				sb.WriteString(abs[index].String())
				sumPrice += abs[index].Price
				sb.WriteString("|\n")
			}
			sb.WriteString("     —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— —— \n")
			sb.WriteString("    [                                           总  价 : ")
			sb.WriteString(common.Float64ToString(sumPrice))
			sb.WriteString("  \t\t  ")
		}
	}
	return sb.String()
}

// OpenBussiness 开始营业
func OpenBussiness() {
	log.InfoLog("开始营业...")
	key := common.GetDate()
	for {
		input := readConsole()
		switch input {
		case "s": //显示数据库
			showDataBase()
		case "d": //删除上一条记录
			deleteLastRecord()
		case "d1": //删除当天记录
			deleteDayRecord()
		case "d2": //删除指定记录
			deleteRecordByBarcode()
		default:
			account(key, input) //记账
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
	log.InfoLog("input:", input[:len(input)-2])
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
			saveAcountBook(key, common.GetRandBarcode(), price)
		} else {
			log.InfoLog("请输入正确的价格：")
		}
	}
}

// saveAcountBook 记录账本
func saveAcountBook(date string, barcode int64, price float64) bool {
	allabs, abs := getAllabsAndNowab() // 获取账簿和当天账本
	goods := make(map[int64]good)      //条形码做键值
	good := good{Name: "默认商品"}
	if err := db.Get(GOODSNAME, &goods); err == nil { //取得商品库
		good = goods[barcode] //取得当前卖出的商品
	}
	ab := accountBook{barcode, good.Name, price, common.GetTime1()} //初始化账目
	abs = append(abs, ab)                                           //加入账本
	allabs[date] = abs                                              //存入账簿
	if db.Put(ACCOUNTBOOKNAME, allabs) {
		log.InfoLog("存入数据库成功!", ab.String())
	}
	return true
}

// showDataBase 输出数据库
func showDataBase() {
	log.InfoLog("开始输出数据库...")
	// iter := db.Database.NewIterator(nil, nil)
	iter := db.Database.NewIterator(util.BytesPrefix([]byte(ACCOUNTBOOKNAME)), nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		// key := iter.Key()
		value := iter.Value()
		var abs map[string][]accountBook
		if err := json.Unmarshal(value, &abs); err != nil {
			db.Database.Delete([]byte(ACCOUNTBOOKNAME), nil)
			log.ErrorLog("反格式化数据错误,value：", abs, ".开始执行删除")
		}
		log.InfoLog(ToString(abs))
	}
	iter.Release()
	// err = iter.Error()
}

// getAllabsAndNowab 获取账簿和当天账本
func getAllabsAndNowab() (map[string][]accountBook, []accountBook) {
	allabs := make(map[string][]accountBook)                 // 日期做键值,每天都是一个小账本
	abs := []accountBook{}                                   // 新建当天账本
	if err := db.Get(ACCOUNTBOOKNAME, &allabs); err == nil { //取得账簿
		abs = allabs[common.GetDate()] //已有当天账本则取值当天账本
	}
	return allabs, abs
}

// deleteRecordByBarcode 根据条形码清除当天记录
func deleteRecordByBarcode() {
	log.InfoLog("请输入要删除的条形码：")
	input := readConsole()
	if barcode, isBarcode := checkBarcode(input); isBarcode {
		allabs, abs := getAllabsAndNowab() // 获取账簿和当天账本
		if len(abs) > 0 {
			var ab accountBook
			var isDelete = false
			for index := 0; index < len(abs); index++ {
				if barcode == abs[index].BarCode {
					ab = abs[index]
					abs = append(abs[:index], abs[index+1:]...)
					isDelete = true
					break
				}
			}
			if isDelete {
				allabs[common.GetDate()] = abs
				if db.Put(ACCOUNTBOOKNAME, allabs) {
					log.InfoLog("删除指定记录成功!", ab.String())
				}
			} else {
				log.InfoLog("找不到数据可以删除!请仔细核对条形码")
			}
		} else {
			log.InfoLog("无数据可以删除!")
		}
	} else {
		log.InfoLog("不正确的条形码,请求拒绝!")
	}

}

// deleteLastRecord 清除上一条记录
func deleteLastRecord() {
	log.InfoLog("确认是否删除上一条记录?(y:是/n:否)")
	if input := readConsole(); input == "y" {
		allabs, abs := getAllabsAndNowab() // 获取账簿和当天账本
		if len(abs) > 0 {
			ab := abs[len(abs)-1]
			abs = abs[:len(abs)-1]
			allabs[common.GetDate()] = abs
			if db.Put(ACCOUNTBOOKNAME, allabs) {
				log.InfoLog("删除上一条记录成功!", ab.String())
			}
		} else {
			log.InfoLog("无数据可以删除!")
		}
	}
}

// deleteLastRecord 清除当天记录
func deleteDayRecord() {
	log.InfoLog("确认是否删除当天记录?(y:是/n:否)")
	if input := readConsole(); input == "y" {
		allabs, _ := getAllabsAndNowab() // 获取账簿和当天账本
		allabs[common.GetDate()] = []accountBook{}
		if db.Put(ACCOUNTBOOKNAME, allabs) {
			log.InfoLog("删除当天记录成功!")
		}
	}
}
