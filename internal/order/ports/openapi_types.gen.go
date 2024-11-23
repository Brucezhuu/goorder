// Package ports provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package ports

// CreateOrderRequest defines model for CreateOrderRequest.
type CreateOrderRequest struct {
	CustomerId string             `json:"customer_id"`
	Items      []ItemWithQuantity `json:"items"`
}

// Error defines model for Error.
type Error struct {
	Message *string `json:"message,omitempty"`
}

// Item defines model for Item.
type Item struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	PriceId  string `json:"price_id"`
	Quantity int32  `json:"quantity"`
}

// ItemWithQuantity defines model for ItemWithQuantity.
type ItemWithQuantity struct {
	Id       string `json:"id"`
	Quantity int32  `json:"quantity"`
}

// Order defines model for Order.
type Order struct {
	CustomerId  string `json:"customer_id"`
	Id          string `json:"id"`
	Items       []Item `json:"items"`
	PaymentLink string `json:"payment_link"`
	Status      string `json:"status"`
}

// Response defines model for Response.
type Response struct {
	Data    map[string]interface{} `json:"data"`
	Errno   int                    `json:"errno"`
	Message string                 `json:"message"`
	TraceId string                 `json:"trace_id"`
}

// PostCustomerCustomerIdOrdersJSONRequestBody defines body for PostCustomerCustomerIdOrders for application/json ContentType.
type PostCustomerCustomerIdOrdersJSONRequestBody = CreateOrderRequest
