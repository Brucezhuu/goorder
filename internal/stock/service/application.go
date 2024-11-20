package service

import (
	"context"

	"github.com/Brucezhuu/goorder/internal/common/metrics"
	"github.com/Brucezhuu/goorder/internal/stock/adapters"
	"github.com/Brucezhuu/goorder/internal/stock/app"
	"github.com/Brucezhuu/goorder/internal/stock/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
