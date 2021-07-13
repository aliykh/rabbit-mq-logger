# rabbit-mq-logger
RabbitMQ logging



make run in log-emitter
above code builds and runs our application which emits logs with types in relation to 1-100 nums
like from 1-30 error logs -> from 31-60 info and the rest debug - smth like that


make build and run with parameters in log-all: i.e. make build && ./bin/app -k='#' for listening all types of logs

the above code builds and runs consumer application in which logs in the rabbitmq queue will be consumed
param -k= '#' for all log types
          'error.*'
          '*.debug.*
          '*.info.*
          
ideally you should run 4 programs with the given params above
