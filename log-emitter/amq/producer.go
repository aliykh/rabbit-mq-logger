package amq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Body struct {
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

type Message struct {
	ContentType string `json:"content_type"`
	Body        Body   `json:"body"`
}

type LogProducer struct {
	conn       *amqp.Connection
	logger     *zap.Logger
	exType     string
	routingKey string
	exName     string
}

func NewLogProducer(conn *amqp.Connection, logger *zap.Logger, exType, routingKey, exName string) *LogProducer {
	return &LogProducer{
		conn:       conn,
		logger:     logger,
		exType:     exType,
		routingKey: routingKey,
		exName:     exName,
	}
}

//func (l LogProducer) InitializeQueue(ch *amqp.Channel, name string) (*amqp.Queue, error) {
//	q, err := ch.QueueDeclare(
//		name,  // name
//		false, // durable
//		false, // delete when unused
//		false, // exclusive
//		false, // no-wait
//		nil,   // arguments
//	)
//
//	if err != nil {
//		return nil, err
//	}
//	return &q, nil
//}

func (l LogProducer) InitializeExchange(ch *amqp.Channel) error {
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

func (l LogProducer) Emit(msg Message) (err error) {

	ch, err := l.conn.Channel()

	if err != nil {
		return
	}

	defer ch.Close()

	err = l.InitializeExchange(ch)
	if err != nil {
		return
	}


	b, err := json.Marshal(msg.Body)
	if err != nil {
		l.logger.Error("error while marshaling body of the message")
		return
	}

	err = ch.Publish(
		l.exName,
		l.routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: msg.ContentType,
			Body:        b,
		})

	return
}
