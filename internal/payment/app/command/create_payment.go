package command

import (
	"context"

	"github.com/Brucezhuu/goorder/internal/common/decorator"
	"github.com/Brucezhuu/goorder/internal/common/genproto/orderpb"
	"github.com/Brucezhuu/goorder/internal/common/metrics"
	"github.com/Brucezhuu/goorder/internal/payment/domain"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGRPC OrderService
}

func NewCreatePaymentHandler(
	processor domain.Processor,
	orderGRPC OrderService,
	logger *logrus.Entry,
	metricsClient metrics.TodoMetrics,
) CreatePaymentHandler {
	return decorator.ApplyCommandDecorators[CreatePayment, string](createPaymentHandler{
		processor: processor,
		orderGRPC: orderGRPC,
	},
		logger,
		metricsClient)
}

func (c createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("create payment link for order: %s success, payment link: %s", cmd.Order.ID, link)
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment",
		Items:       cmd.Order.Items,
		PaymentLink: link,
	}
	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
	return link, err
}
