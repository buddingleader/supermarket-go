package good

import (
	"errors"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wangff15386/supermarket-go/local/mongodb"
	"github.com/wangff15386/supermarket-go/log"
	"gopkg.in/mgo.v2/bson"
)

// collection name
const (
	NAME  = "GOODs"
	COUNT = "GOOD_COUNT"
)

// Service for ACID good
type Service struct {
	logger *logrus.Entry
}

// NewGoodService for new good service
func NewGoodService() *Service {
	return &Service{
		logger: log.GetLogger("good"),
	}
}

// Good standard good
type Good struct {
	Barcode           string  `bson:"barcode"`
	Name              string  `bson:"name"`
	Specification     string  `bson:"specification"`
	Unit              string  `bson:"unit"`
	LastPurchasePrice float64 `bson:"lpprice"`
	MaximumSalePrice  float64 `bson:"msprice"`
	Quantity          string  `bson:"quantity"`
	Stock             uint    `bson:"stock"`
}

// GetGood get good from mongo
func (gs *Service) GetGood(barcode string) (*Good, error) {
	database, ctx, cancel := mongodb.GetConn()
	defer cancel()

	result, good := database.Collection(NAME).FindOne(ctx, &bson.M{"barcode": barcode}), &Good{}
	if err := result.Decode(good); err != nil {
		if strings.Contains(err.Error(), mongodb.FOUNTNOT) {
			return nil, nil
		}
		return nil, err
	}
	return good, nil
}

// GetGoods get all the goods from mongo
func (gs *Service) GetGoods() ([]*Good, error) {
	database, ctx, cancel := mongodb.GetConn()
	defer cancel()

	cur, err := database.Collection(NAME).Find(ctx, nil)
	defer cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	goods, good := make([]*Good, 0), &Good{}
	for cur.Next(ctx) {
		err := cur.Decode(good)
		if err != nil {
			return nil, err
		}

		goods = append(goods, good)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return goods, nil
}

// PutGood put a good to mongo
func (gs *Service) PutGood(good *Good) error {
	database, ctx, cancel := mongodb.GetConn()
	defer cancel()

	g, err := gs.GetGood(good.Barcode)
	if err != nil {
		return err
	}
	if g != nil {
		if _, err = database.Collection(NAME).ReplaceOne(ctx, g, good); err != nil {
			return err
		}
		gs.logger.Infof("repalce a good[%s] to good[%s]", g, good)
	} else {
		if _, err = database.Collection(NAME).InsertOne(ctx, good); err != nil {
			return err
		}

		gs.logger.Infof("insert a good[%s]", good)
	}

	return nil
}

// DeleteGood delete a good
func (gs *Service) DeleteGood(barcode string) error {
	database, ctx, cancel := mongodb.GetConn()
	defer cancel()

	good, err := gs.GetGood(barcode)
	if err != nil {
		return err
	}

	if good == nil {
		return errors.New("Cannnot found the Good")
	}

	if _, err = database.Collection(NAME).DeleteOne(ctx, good); err != nil {
		return err
	}
	gs.logger.Infof("delete a good[%s]", good)

	return nil
}

// GetGoodsCount get goods count
func (gs *Service) GetGoodsCount() (int64, error) {
	database, ctx, cancel := mongodb.GetConn()
	defer cancel()

	count, err := database.Collection(NAME).Count(ctx, nil)
	if err != nil {
		return 0, err
	}
	return count, nil
}
