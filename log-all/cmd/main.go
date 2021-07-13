package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aliykh/rabbitmq-logger/log-all/amq"
	"github.com/aliykh/rabbitmq-logger/log-all/controller"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

func getLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	logger.Info("logger instantiated.....")
	return logger
}

func main() {

	//Listen os interruption for graceful shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		osCall := <-ch
		log.Printf("OS Interruption Call: %v\n", osCall)
		cancel()
	}()

	//Get routing key
	if len(os.Args) < 2 || os.Args[1] == "" {
		log.Printf("%v\n", os.Args)
		log.Fatalln("please provide the program with routing key for the topic [rabbit MQ] i.e. app -k=error.* ")
	}

	routingKey := flag.String("k", "error.*", "set route key for the logger")
	flag.Parse()

	//Set env configurations
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
	exName := os.Getenv("exchange_name")
	exType := os.Getenv("exchange_type")

	//Dial rabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, psw, host, port))
	if err != nil {
		logger.Error("rabbitmq err", zap.Any("any", err.Error()))
		return
	}
	defer conn.Close()
	logger.Info("rabbit-MQ connected...")

	//APP logic
	c := amq.NewLogConsumer(logger, *routingKey, exName, exType)

	ctrl := controller.NewLogController(logger, c, conn)

	go ctrl.ReceiveLogs(ctx)

	//go func() {
	//	<-ctx.Done()
	//	fmt.Println("Prepare for graceful shutdown...")
	//	time.Sleep(time.Second * 4)
	//	fmt.Println("SHUTDOWN.....")
	//
	//}()
	<-ctx.Done()
	fmt.Println("it is done...")

}
