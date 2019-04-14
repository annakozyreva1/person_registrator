package main

import (
	"log"
	"github.com/streadway/amqp"
	"reflect"
	"time"
	"fmt"
	"net"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//"amqp://gpiicshf:Ja2mQPOi2Mz25K7dmJtwpZOGfu-WaH3v@gopher.rmq.cloudamqp.com/gpiicshf

func main() {
	conn, err := amqp.Dial("amqp://gpiicshf:Ja2mQPOi2Mz25K7dmJtwpZOGfu-WaH3v@gopher.rmq.cloudamqp.com/gpiicshf")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	/*q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)*/
	failOnError(err, "Failed to declare a queue")

	body := "hello"
	r := make(chan *amqp.Error)
	ch.NotifyClose(r)
	go func() {
		select {
		case rr :=<- r:
			{
				fmt.Printf("%v", rr)
			}
		}
	}()
	for {
		err = ch.Publish(
			"",     // exchange
			"hjhj", // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		//failOnError(err, "Failed to publish a message")
		if err != nil {
			log.Printf(" [x] Sent %s %+v", body, reflect.TypeOf(err))
			if _, ok := err.(*net.OpError); ok {

			}
		}

		time.Sleep(time.Second)
	}
}

