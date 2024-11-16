package query

import (
	"context"
	"github.com/Brucezhuu/goorder/internal/common/genproto/orderpb"
	"github.com/Brucezhuu/goorder/internal/common/genproto/stockpb"
)

type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error)
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
