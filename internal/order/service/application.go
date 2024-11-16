package service

import (
	"context"
	grpcClient "github.com/Brucezhuu/goorder/internal/common/client"
	"github.com/Brucezhuu/goorder/internal/common/metrics"
	"github.com/Brucezhuu/goorder/internal/order/adapters"
	"github.com/Brucezhuu/goorder/internal/order/adapters/grpc"
	"github.com/Brucezhuu/goorder/internal/order/app"
	"github.com/Brucezhuu/goorder/internal/order/app/command"
	"github.com/Brucezhuu/goorder/internal/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
