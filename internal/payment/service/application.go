package service

import (
	"context"
	grpcClient "github.com/Brucezhuu/goorder/internal/common/client"
	"github.com/Brucezhuu/goorder/internal/common/metrics"
	"github.com/Brucezhuu/goorder/internal/payment/adapters"
	"github.com/Brucezhuu/goorder/internal/payment/app"
	"github.com/Brucezhuu/goorder/internal/payment/app/command"
	"github.com/Brucezhuu/goorder/internal/payment/domain"
	"github.com/Brucezhuu/goorder/internal/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	memoryProcessor := processor.NewInmemProcessor()
	return newApplication(ctx, orderGRPC, memoryProcessor), func() {
		_ = closeOrderClient()
	}
}
func newApplication(ctx context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	// 这里函数的传参都依赖于接口
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}
