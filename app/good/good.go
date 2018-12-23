package good

import (
	"github.com/kataras/iris"
	"github.com/wangff15386/supermarket-go/server"
)

// FOUNDNOT cannot found the good
const FOUNDNOT = "Cannot found the good! \n商品不存在！"

// Get scan barcode and get the good information
func Get(ctx iris.Context) {
	barcode := ctx.Params().GetStringDefault("barcode", "")
	if barcode == "" {
		ctx.WriteString(FOUNDNOT)
		return
	}

	gs := server.NewGoodService()
	g, err := gs.GetGood(barcode)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	if g == nil {
		ctx.WriteString(FOUNDNOT)
		return
	}
	ctx.JSON(g)
}

// GetAll get all the goods information
func GetAll(ctx iris.Context) {
	gs := server.NewGoodService()
	g, err := gs.GetGoods()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.JSON(g)
}
