package controller

import (
	"fmt"
	"github.com/aliykh/rabbitmq-logger/log-emitter/amq"
	"go.uber.org/zap"
	"strconv"
)

type logController struct {
	logger    *zap.Logger
	producers []*amq.LogProducer
}

func NewLogController(logger *zap.Logger, producers []*amq.LogProducer) *logController {
	return &logController{logger: logger, producers: producers}
}

func (l *logController) EmitLogs() {

	for i := 1; i <= 100; i++ {
		if i <= 30 {
			err := l.producers[0].Emit(amq.Message{
				ContentType: "application/json",
				Body: amq.Body{
					Type: "error",
					Msg: strconv.Itoa(i),
				},
			})
			if err != nil {
				l.logger.Error(fmt.Sprintf("error while emitting the message %s", err.Error()))
			}
		} else if i > 30 && i <= 60 {
			err := l.producers[1].Emit(amq.Message{
				ContentType: "application/json",
				Body: amq.Body{
					Type: "info",
					Msg: strconv.Itoa(i),
				},
			})
			if err != nil {
				l.logger.Error(fmt.Sprintf("error while emitting the message %s", err.Error()))
			}
		} else {
			err := l.producers[2].Emit(amq.Message{
				ContentType: "application/json",
				Body: amq.Body{
					Type: "debug",
					Msg: strconv.Itoa(i),
				},
			})
			if err != nil {
				l.logger.Error(fmt.Sprintf("error while emitting the message %s", err.Error()))
			}
		}
	}

	//for _, v := range l.producers {
	//	id := uuid.NewV4()
	//	err := v.Emit(amq.Message{
	//		ContentType: "application/json",
	//		Body: amq.Body{
	//			Msg: id.String(),
	//		},
	//	})
	//
	//	if err != nil {
	//		l.logger.Error(fmt.Sprintf("error while emitting the message %s", err.Error()))
	//	}
	//}

}
