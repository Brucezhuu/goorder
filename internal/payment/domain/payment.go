package domain

import (
	"context"
	"github.com/Brucezhuu/goorder/internal/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
