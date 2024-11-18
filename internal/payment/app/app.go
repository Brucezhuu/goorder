package app

import "github.com/Brucezhuu/goorder/internal/payment/app/command"

type Application struct {
	Commands Commands
}

type Commands struct {
	CreatePayment command.CreatePaymentHandler
}
