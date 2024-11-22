package main

import (
	"context"
	"github.com/Brucezhuu/goorder/internal/common/tracing"

	"github.com/Brucezhuu/goorder/internal/common/broker"
	_ "github.com/Brucezhuu/goorder/internal/common/config"
	"github.com/Brucezhuu/goorder/internal/common/logging"
	"github.com/Brucezhuu/goorder/internal/common/server"
	"github.com/Brucezhuu/goorder/internal/payment/infrastructure/consumer"
	"github.com/Brucezhuu/goorder/internal/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
}

func main() {
	serviceName := viper.GetString("payment.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	serverType := viper.GetString("payment.server-to-run")

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)
	switch serverType {
	case "http":
		server.RunHTTPServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported server type: grpc")
	default:
		logrus.Panic("unreachable code")
	}
}
