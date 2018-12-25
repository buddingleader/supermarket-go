package good

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangff15386/supermarket-go/common/utils"
	"github.com/wangff15386/supermarket-go/conf"
)

func init() {
	conf.Initial("../../conf/app.conf")
}

func Test_PutGood(t *testing.T) {
	gs, good := NewGoodService(), newTestGood()
	err := gs.PutGood(good)
	assert.NoError(t, err)

	good.MaximumSalePrice = 88
	err = gs.PutGood(good)
	assert.NoError(t, err)
}

func Test_GetGood(t *testing.T) {
	gs, good := NewGoodService(), newTestGood()
	err := gs.PutGood(good)
	assert.NoError(t, err)

	// Correct
	g, err := gs.GetGood(good.Barcode)
	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, good.MaximumSalePrice, g.MaximumSalePrice)

	// Invalid barcode
	g1, err := gs.GetGood("agdfbsg")
	assert.NoError(t, err)
	assert.Nil(t, g1)
}

func Test_GetGoods(t *testing.T) {
	gs := NewGoodService()

	goods, err := gs.GetGoods()
	assert.NoError(t, err)
	count, err := gs.GetGoodsCount()
	assert.NoError(t, err)
	assert.Equal(t, count, int64(len(goods)))
}

func Test_DeleteGood(t *testing.T) {
	gs, good := NewGoodService(), newTestGood()
	err := gs.PutGood(good)
	assert.NoError(t, err)

	err = gs.DeleteGood(good.Barcode)
	assert.NoError(t, err)
}

func newTestGood() *Good {
	return &Good{
		Barcode:          utils.GetRandBarcodeS(),
		Name:             "采果姑娘（鲜园坊）",
		MaximumSalePrice: 6.5,
		Quantity:         "250ml",
		Specification:    "盒",
	}
}
