package command

import (
	"context"
	"github.com/Brucezhuu/goorder/internal/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
