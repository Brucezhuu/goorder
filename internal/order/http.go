package main

import (
	"fmt"
	"github.com/Brucezhuu/goorder/internal/common"
	client "github.com/Brucezhuu/goorder/internal/common/client/order"
	"github.com/Brucezhuu/goorder/internal/order/app"
	"github.com/Brucezhuu/goorder/internal/order/app/command"
	"github.com/Brucezhuu/goorder/internal/order/app/query"
	"github.com/Brucezhuu/goorder/internal/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var (
		req  client.CreateOrderRequest
		err  error
		resp struct {
			CustomerID  string `json:"customer_id"`
			OrderID     string `json:"order_id"`
			RedirectURL string `json:"redirect_url"`
		}
	)
	defer func() {
		H.Response(c, err, &resp)
	}()
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp.CustomerID = req.CustomerID
	resp.RedirectURL = fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID)
	resp.OrderID = r.OrderID
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp interface{}
	)
	defer func() {
		H.Response(c, err, resp)
	}()
	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{CustomerID: customerID, OrderID: orderID})
	if err != nil {
		return
	}
	resp = convertor.NewOrderConvertor().EntityToClient(o)
}
