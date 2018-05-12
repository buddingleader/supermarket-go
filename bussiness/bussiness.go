package bussiness

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"supermarket-go/bussiness/good"
	"supermarket-go/common/utils"
	db "supermarket-go/local/leveldb"
	"supermarket-go/log"

	"github.com/syndtr/goleveldb/leveldb/util"
)

// 对象键值
const (
	ACCOUNTBOOKNAME = "accountBook" //账簿名称
	GOODSNAME       = "goods"       //商品库名称
)

var (
	cacheAccount []accountBook
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
func absToString(date string, abs []accountBook) string {
	sb := bytes.Buffer{}
	var sumPrice float64
	if len(abs) > 0 {
		sumPrice = abs[0].Price
		sb.WriteString("<--日期：")
		sb.WriteString(date)
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
		sb.WriteString("                                              总  价 : ")
		sb.WriteString(common.Float64ToString(sumPrice))
		sb.WriteString("  \t\t  \n")
	}
	return sb.String()
}

// ToString 账簿输出
func ToString(allabs map[string][]accountBook) string {
	sb := bytes.Buffer{}
	sb.WriteString(string(ACCOUNTBOOKNAME))
	sb.WriteString(" 账簿：]\n")
	// var sumPrice float64
	for i, abs := range allabs {
		sb.WriteString(absToString(i, abs))
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
		case "s": //显示当天数据库
			showNowDataBase()
		case "s1": //显示全部数据库
			showDataBase()
		case "s2": //显示指定日期的数据库
			showTheDayDataBase()
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
		accountBarCode(key, barcode)
	} else {
		log.InfoLog("非条形码输入，尝试匹配价格...")
		if price, isPrice := checkPrice(input); isPrice {
			saveAcountBook(key, common.GetRandBarcode(), price, true)
		} else if input == "" {
			for _, ab := range cacheAccount {
				saveAcountBook(key, ab.BarCode, ab.Price, false)
			}
			fmt.Println("本次记账清单：", absToString(key, cacheAccount))
			cacheAccount = []accountBook{}
		} else {
			log.InfoLog("请输入正确的价格：")
		}
	}
}

func accountBarCode(key string, barcode int64) {
	log.InfoLog("标准条形码输入，请输入价格：")
	msp := good.SellPrice{}
	g, err := good.GetGoodPrice(barcode)
	if err == nil {
		fmt.Println("历史出售价格       历史出售时间           价格出售次数")
		for _, sp := range g.OutPrice {
			if sp.Count > msp.Count {
				msp = sp
			}
			fmt.Printf("  %.2f\t\t%s\t\t%d\n", sp.Price, sp.Time, sp.Count)
		}
		ab := accountBook{
			BarCode: barcode,
			Name:    g.Name,
			Price:   msp.Price,
		}
		cacheAccount = append(cacheAccount, ab)
		fmt.Println("缓存的账本：", absToString(key, cacheAccount))
		log.InfoLog("已使用推荐价格：", msp.Price, " 你也可以重新输入价格或继续输入条形码或按Enter提交记账：")
		input1 := readConsole()
		if barcode, isBarcode := checkBarcode(input1); isBarcode {
			accountBarCode(key, barcode)
		} else if price, isPrice := checkPrice(input1); isPrice {
			ab.Price = price
			cacheAccount = append(cacheAccount[:len(cacheAccount)-1], ab)
			fmt.Println("缓存的账本：", absToString(key, cacheAccount), " 你可以继续输入条形码或按Enter提交记账：")
		} else if input1 == "" {
			for _, ab := range cacheAccount {
				saveAcountBook(key, ab.BarCode, ab.Price, false)
			}
			fmt.Println("本次记账清单：", absToString(key, cacheAccount))
			cacheAccount = []accountBook{}
		}
	}
}

// saveAcountBook 记录账本
func saveAcountBook(date string, barcode int64, price float64, isRandom bool) bool {
	allabs, abs := getAllabsAndNowab() // 获取账簿和当天账本
	g, err := good.GetGood(barcode)
	if err != nil {
		log.ErrorLog("获取商品失败!", err)
		return false
	}
	ab := accountBook{barcode, g.Name, price, common.GetTime1()} //初始化账目
	abs = append(abs, ab)                                        //加入账本
	allabs[date] = abs                                           //存入账簿
	if db.Put(ACCOUNTBOOKNAME, allabs) {
		log.InfoLog("存入数据库成功!", ab.String())
		if !isRandom {
			good.PutGoodPrice(g, price)
		}
	}
	return true
}

// showNowDataBase 输出当天数据库
func showNowDataBase() {
	log.InfoLog("开始输出数据库...")
	// iter := db.Database.NewIterator(nil, nil)
	iter := db.Database.NewIterator(util.BytesPrefix([]byte(ACCOUNTBOOKNAME)), nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		// key := iter.Key()
		value := iter.Value()
		var allabs map[string][]accountBook
		if err := json.Unmarshal(value, &allabs); err != nil {
			db.Database.Delete([]byte(ACCOUNTBOOKNAME), nil)
			log.ErrorLog("反格式化数据错误,value：", allabs, ".开始执行删除")
		}
		log.InfoLog(absToString(common.GetDate(), allabs[common.GetDate()]))
	}
	iter.Release()
	// err = iter.Error()
}

// showDataBase 输出全部数据库
func showDataBase() {
	log.InfoLog("开始输出数据库...")
	// iter := db.Database.NewIterator(nil, nil)
	iter := db.Database.NewIterator(util.BytesPrefix([]byte(ACCOUNTBOOKNAME)), nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		// key := iter.Key()
		value := iter.Value()
		var allabs map[string][]accountBook
		if err := json.Unmarshal(value, &allabs); err != nil {
			db.Database.Delete([]byte(ACCOUNTBOOKNAME), nil)
			log.ErrorLog("反格式化数据错误,value：", allabs, ".开始执行删除")
		}
		log.InfoLog(ToString(allabs))
	}
	iter.Release()
	// err = iter.Error()
}

// showTheDayDataBase() 输出制定日期数据库
func showTheDayDataBase() {
	log.InfoLog("请输入日期[格式：20180428]查询数据库：")
	input := readConsole()
	allabs, _ := getAllabsAndNowab()
	abs := allabs[input]
	log.InfoLog(absToString(input, abs))
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
