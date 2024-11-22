package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel"

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
			c.handleMessage(ch, msg, q)
		}
	}()
	<-forever // 在这里永远阻塞住
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("failed to unmarshall msg to order, err=%v", err)
		return
	}
	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		logrus.Infof("failed to create payment, err=%v", err)
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
		}
		return
	}

	span.AddEvent("payment.created")
	logrus.Info("consume success")
}
