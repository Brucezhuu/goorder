package main

import (
	"log"

	"github.com/Brucezhuu/goorder/internal/common/config"
	"github.com/Brucezhuu/goorder/internal/common/genproto/orderpb"
	"github.com/Brucezhuu/goorder/internal/common/server"
	"github.com/Brucezhuu/goorder/internal/order/ports"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")
	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer()
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
