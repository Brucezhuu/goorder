package main

import (
	"fmt"
	"github.com/Brucezhuu/goorder/internal/common/config"
	"github.com/Brucezhuu/goorder/internal/common/genproto/stockpb"
	"github.com/Brucezhuu/goorder/internal/common/server"
	"github.com/Brucezhuu/goorder/internal/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}
func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")
	fmt.Println(serviceName, serverType)
	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer()
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		// 暂时不用
	default:
		panic("unexpected server type")
	}
}
