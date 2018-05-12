package good

import (
	"fmt"
	"log"
	"testing"
)

func Test_putExcel(t *testing.T) {
	goods, err := GetGoods()
	if err != nil {
		log.Fatal(err)
	}
	ShowGoods(goods)

	path := "good_test.xlsx"
	if b, err := putExcel(path); !b {
		log.Fatal(err)
	}

	goods, err = GetGoods()
	if err != nil {
		log.Fatal(err)
	}
	ShowGoods(goods)
}

func Test_ProcessGood(t *testing.T) {
	good, err := GetGood(6914973604469)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(good)

	PutGoodPrice(good, 14)

	good, err = GetGood(6914973604469)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(good)

	fmt.Println("good.OutPrice", good.OutPrice)
	SortSellPrice(good.OutPrice, func(p, q *SellPrice) bool {
		return p.Time > q.Time
	})
	fmt.Println("good.OutPrice", good.OutPrice)

	good, err = GetGood(6914973604467)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(good)

}
