package good

import (
	"errors"
	"fmt"
	"sort"
	common "supermarket-go/common/utils"
	db "supermarket-go/local/leveldb"
	"supermarket-go/log"

	"github.com/tealeg/xlsx"
)

// 对象键值
const (
	GOODSNAME     = "goods" //商品库名称
	EXCELFILEPATH = "GoodsBarcode_excel.xlsx"
)

// Good 商品
type Good struct {
	Barcode       int64       `json:"barCode"`
	Name          string      `json:"Name"`
	InPrice       float64     `json:"InPrice"`
	OutPrice      []SellPrice `json:"OutPrice"`
	Quantity      string      `json:"Quantity"`
	Specification string      `json:"Specification"`
}

func (g *Good) String() string {
	return fmt.Sprintf("Good[Barcode=%d, Name=%s, InPrice=%.2f, OutPrice=%v, Quantity=%s, Specification=%s]", g.Barcode, g.Name, g.InPrice, g.OutPrice, g.Quantity, g.Specification)
}

// SellPrice 销售价格
type SellPrice struct {
	Price float64
	Time  string
	// 销售次数
	Count int64
}

func (sp *SellPrice) String() string {
	return fmt.Sprintf("SellPrice[Price=%.2f, Time=%s, Count=%d]", sp.Price, sp.Time, sp.Count)
}

// SellPriceWrapper 排序辅助结构
type SellPriceWrapper struct {
	prices []SellPrice
	by     func(p, q *SellPrice) bool
}

// Len 重写用于排序的Len方法
func (spw SellPriceWrapper) Len() int { // 重写 Len() 方法
	return len(spw.prices)
}

// Swap 重写用于排序的Swap方法
func (spw SellPriceWrapper) Swap(i, j int) { // 重写 Swap() 方法
	spw.prices[i], spw.prices[j] = spw.prices[j], spw.prices[i]
}

// Less 重写用于排序的Less方法
func (spw SellPriceWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return spw.by(&spw.prices[i], &spw.prices[j])
}

// SortSellPrice 封装成 SortPerson 方法
func SortSellPrice(prices []SellPrice, by func(p, q *SellPrice) bool) {
	sort.Sort(SellPriceWrapper{prices, by})
}

// GetGoods 从数据库中取得商品库
func GetGoods() (map[int64]Good, error) {
	goods := make(map[int64]Good)                     //条形码做键值
	if err := db.Get(GOODSNAME, &goods); err != nil { //取得商品库
		log.InfoLog(err, "初始化Goods")
		PutGoods(goods)
		return goods, err
	}
	return goods, nil
}

// PutGoods 存储商品库到数据库中
func PutGoods(goods map[int64]Good) bool {
	return db.Put(GOODSNAME, goods)
}

func putExcel(path string) (bool, error) {
	goods, err := GetGoods()
	if err != nil { //取得商品库
		log.ErrorLog(err)
		return false, err
	}
	if path == "" {
		path = EXCELFILEPATH
	}
	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
		return false, err
	}
	for _, sheet := range xlFile.Sheets {
		fmt.Printf("Sheet Name: %s\n", sheet.Name)
		for _, row := range sheet.Rows {
			barcode, err := row.Cells[0].Int64()
			if err != nil {
				log.ErrorLog(err)
				continue
			}
			if _, ok := goods[barcode]; !ok {
				good := Good{
					Barcode:       barcode,
					Name:          row.Cells[1].String(),
					Quantity:      row.Cells[2].String(),
					Specification: row.Cells[3].String(),
				}
				goods[barcode] = good
			}
		}
	}
	return PutGoods(goods), nil
}

// ShowGoods 展示商品库
func ShowGoods(goods map[int64]Good) {
	for _, good := range goods {
		fmt.Println("good:", good.String())
	}
}

// GetGood 从商品库中取得商品
func GetGood(barcode int64) (Good, error) {
	goods, err := GetGoods()
	if err != nil {
		log.ErrorLog(err)
		return Good{}, err
	}
	good, ok := goods[barcode]
	if !ok {
		good = Good{
			Barcode: barcode,
			Name:    "未知商品",
		}
	}
	return good, nil
}

// PutGood 存储商品到商品库中，并更新销售指导价
func PutGood(good Good) (bool, error) {
	goods, err := GetGoods()
	if err != nil {
		log.ErrorLog(err)
		return false, err
	}
	goods[good.Barcode] = good
	return PutGoods(goods), nil
}

// GetGoodPrice 从商品库中取得商品的历史销售价格，并按时间倒序排序
func GetGoodPrice(barcode int64) (Good, error) {
	good, err := GetGood(barcode)
	if err != nil {
		log.ErrorLog(err)
		return Good{}, err
	}
	SortSellPrice(good.OutPrice, func(p, q *SellPrice) bool {
		return p.Time > q.Time
	})
	return good, nil
}

// PutGoodPrice 存储商品到商品库中，并更新销售指导价
func PutGoodPrice(good Good, price float64) (bool, error) {
	if price <= 0 {
		return false, errors.New("金额小于等于0")
	}

	prices := good.OutPrice
	for index, sp := range prices {
		if sp.Price == price {
			sp.Count++
			sp.Time = common.GetTimeStamp1()
			prices[index] = sp
			good.OutPrice = prices
			return PutGood(good)
		}
	}

	good.OutPrice = append(prices, SellPrice{
		Price: price,
		Time:  common.GetTimeStamp1(),
		Count: 1,
	})
	return PutGood(good)
}

// DelGoodPrice 删除商品库中商品的销售指导价
func DelGoodPrice(barCode int64, price float64) (bool, error) {
	good, err := GetGood(barCode)
	if err != nil {
		log.ErrorLog(err)
		return false, err
	}
	prices := good.OutPrice
	for index, sp := range prices {
		if sp.Price == price {
			good.OutPrice = append(prices[:index], prices[index+1:]...)
			return PutGood(good)
		}
	}
	return false, nil
}
