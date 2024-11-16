package app

import (
	"github.com/Brucezhuu/goorder/internal/order/app/command"
	"github.com/Brucezhuu/goorder/internal/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
