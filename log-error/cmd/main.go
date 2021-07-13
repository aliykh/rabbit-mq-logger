package main

import (
	"fmt"
	"github.com/aliykh/rabbitmq-logger/log-error/amq"
	"github.com/aliykh/rabbitmq-logger/log-error/controller"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"os"
)

func getLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	logger.Info("logger instantiated.....")
	return logger
}

func main() {

	err := godotenv.Load()

	if err != nil {
		return
	}

	logger := getLogger()
	defer logger.Sync()

	host := os.Getenv("rabbitmq_host")
	port := os.Getenv("rabbitmq_port")
	user := os.Getenv("rabbitmq_user")
	psw := os.Getenv("rabbitmq_password")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, psw, host, port))
	if err != nil {
		logger.Error("rabbitmq err", zap.Any("any", err.Error()))
		return
	}
	defer conn.Close()
	logger.Info("rabbit-MQ connected...")

	c := amq.NewLogConsumer(logger, "error.*", "logs", "topic")

	ctrl := controller.NewLogController(logger, c, conn)

	ctrl.ReceiveLogs()

}
