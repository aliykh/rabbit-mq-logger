package controller

import (
	"fmt"
	amq "github.com/aliykh/rabbitmq-logger/log-error/amq"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type logController struct {
	logger    *zap.Logger
	consumer *amq.LogConsumer
	conn   *amqp.Connection
}


func NewLogController(logger *zap.Logger, consumer *amq.LogConsumer, conn *amqp.Connection) *logController {
	return &logController{logger: logger, consumer: consumer, conn: conn}
}

func (l logController) ReceiveLogs() {

	ch, err := l.conn.Channel()

	if err != nil {
		l.logger.Error(fmt.Sprintf("Error while creating channel %s", err.Error()))
		return
	}

	defer ch.Close()


	err = l.consumer.InitializeExchange(ch)
	if err != nil {
		l.logger.Error(fmt.Sprintf("Error while initializing exchange %s", err.Error()))
		return
	}
	
	err = l.consumer.Receive(ch)
	if err != nil {
		l.logger.Error(fmt.Sprintf("Error while receiving message %s", err.Error()))
		return
	}


}