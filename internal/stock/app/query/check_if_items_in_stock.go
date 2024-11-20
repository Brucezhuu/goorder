package query

import (
	"context"

	"github.com/Brucezhuu/goorder/internal/common/decorator"
	"github.com/Brucezhuu/goorder/internal/common/genproto/orderpb"
	domain "github.com/Brucezhuu/goorder/internal/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("nil stockRepo")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*orderpb.Item](
		checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricClient,
	)
}

// TODO: 删掉
var stub = map[string]string{
	"1": "price_1QMWkkKYFdaSkXLdXGwLx3Bp",
	"2": "price_1QM8dOKYFdaSkXLdt5DjTs5A",
}

func (h checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var res []*orderpb.Item
	for _, i := range query.Items {
		// TODO: 改成从数据库 or stripe 获取
		priceID, ok := stub[i.ID]
		if !ok {
			priceID = stub["1"]
		}
		res = append(res, &orderpb.Item{
			ID:       i.ID,
			Quantity: i.Quantity,
			PriceID:  priceID,
		})
	}
	return res, nil
}
