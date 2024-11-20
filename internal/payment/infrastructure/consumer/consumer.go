package consumer

import (
	"context"
	"encoding/json"

	"github.com/Brucezhuu/goorder/internal/common/broker"
	"github.com/Brucezhuu/goorder/internal/common/genproto/orderpb"
	"github.com/Brucezhuu/goorder/internal/payment/app"
	"github.com/Brucezhuu/goorder/internal/payment/app/command"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(application app.Application) *Consumer {
	return &Consumer{
		app: application,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("Fail to consume: queue: %s, err: %v", q.Name, err)
	}
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q)
		}
	}()
	<-forever // 在这里永远阻塞住
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("failed to unmarshall msg to order, err=%v", err)
		_ = msg.Nack(false, false)
		return
	}
	if _, err := c.app.Commands.CreatePayment.Handle(context.TODO(), command.CreatePayment{Order: o}); err != nil {
		// TODO: retry
		logrus.Infof("failed to reate order, err=%v", err)
		_ = msg.Nack(false, false)
		return
	}

	_ = msg.Ack(false)
	logrus.Info("consume success")
}
