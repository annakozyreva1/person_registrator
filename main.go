package main

import (
	"flag"
	"github.com/annakozyreva1/person_registrator/bus"
	"github.com/annakozyreva1/person_registrator/log"
	"github.com/annakozyreva1/person_registrator/person"
	"github.com/annakozyreva1/person_registrator/web"
)

var logger = log.Logger

var (
	busURL     = flag.String("amqp", "", "amqp url as amqp://user:psw@host")
	queue      = flag.String("queue", "person-reg", "bus queue in default exchange")
	webAddress = flag.String("addr", ":7878", "web api address")
)

func main() {
	flag.Parse()
	if *busURL == "" {
		logger.Fatal("need set amqp")
	}
	bus := bus.New(*busURL)
	err := bus.CreateQueue(*queue, true, false)
	if err != nil {
		logger.Fatalf("failed to create queue: %s", err.Error())
	}
	registrator := person.New(bus, *queue)
	web.Run(*webAddress, registrator)
}
