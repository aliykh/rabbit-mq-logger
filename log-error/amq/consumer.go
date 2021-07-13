package amq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
)

type Body struct {
	Msg  string `json:"msg"`
	Type string `json:"type"`
}


type LogConsumer struct {
	logger     *zap.Logger
	routingKey string
	exName     string
	exType     string
}

func NewLogConsumer(logger *zap.Logger, routingKey, exName, exType string) *LogConsumer {
	return &LogConsumer{logger: logger, exName: exName, routingKey: routingKey, exType: exType}
}

func (l LogConsumer) InitializeExchange(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		l.exName,
		l.exType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // argument
	)

	if err != nil {
		return err
	}

	return nil
}

func (l LogConsumer) Receive(ch *amqp.Channel) error {

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name, // queue name
		l.routingKey,     // routing key
		l.exName, // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	delivery, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {

		for v := range delivery {
			data := &Body{}
			err = json.Unmarshal(v.Body, data)
			if err != nil {
				l.logger.Sugar().Errorf("Error while parsing the body %v", err.Error())
			}
			l.logger.Sugar().Errorf("Log: %s; Type: %s\n", data.Msg, data.Type)
		}

	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

	return nil
}
