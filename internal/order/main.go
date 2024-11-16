package main

import (
	"context"
	"github.com/Brucezhuu/goorder/internal/common/config"
	"github.com/Brucezhuu/goorder/internal/common/discovery"
	"github.com/Brucezhuu/goorder/internal/common/genproto/orderpb"
	"github.com/Brucezhuu/goorder/internal/common/server"
	"github.com/Brucezhuu/goorder/internal/order/ports"
	"github.com/Brucezhuu/goorder/internal/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application, cleanup := service.NewApplication(ctx)
	defer cleanup()
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
