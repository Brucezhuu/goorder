package main

import (
	"fmt"
	"github.com/Brucezhuu/goorder/internal/common"
	client "github.com/Brucezhuu/goorder/internal/common/client/order"
	"github.com/Brucezhuu/goorder/internal/common/consts"
	"github.com/Brucezhuu/goorder/internal/common/handler/errors"
	"github.com/Brucezhuu/goorder/internal/order/app"
	"github.com/Brucezhuu/goorder/internal/order/app/command"
	"github.com/Brucezhuu/goorder/internal/order/app/dto"
	"github.com/Brucezhuu/goorder/internal/order/app/query"
	"github.com/Brucezhuu/goorder/internal/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {
	var (
		req  client.CreateOrderRequest
		err  error
		resp dto.CreateOrderResponse
	)
	defer func() {
		H.Response(c, err, &resp)
	}()
	if err := c.ShouldBindJSON(&req); err != nil {
		err = errors.NewWithError(consts.ErrnoBindRequestsError, err)
		return
	}
	if err = H.validate(req); err != nil {
		err = errors.NewWithError(consts.ErrnoRequestsValidateError, err)
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		//err = errors.NewWithError()
		return
	}
	resp = dto.CreateOrderResponse{
		OrderID:     r.OrderID,
		CustomerID:  req.CustomerId,
		RedirectURL: fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}
}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerID string, orderID string) {
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

func (H HTTPServer) validate(req client.CreateOrderRequest) error {
	for _, v := range req.Items {
		if v.Quantity <= 0 {
			return fmt.Errorf("Quantity must be positive, got %d from %s", v.Quantity, v.Id)
		}
	}
	return nil
}
