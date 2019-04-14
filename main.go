package main

import (
	"github.com/annakozyreva1/person_registrator/bus"
	"github.com/annakozyreva1/person_registrator/log"
	"sync"
	"math/rand"
	"sync/atomic"
	"time"
	"github.com/annakozyreva1/person_registrator/web"
)

var logger = log.Logger

//"amqp://gpiicshf:Ja2mQPOi2Mz25K7dmJtwpZOGfu-WaH3v@gopher.rmq.cloudamqp.com/gpiicshf

func main() {
	b := bus.New("amqp://gpiicshf:Ja2mQPOi2Mz25K7dmJtwpZOGfu-WaH3v@gopher.rmq.cloudamqp.com/gpiicshf")
	err := b.CreateQueue("test1", true, false)
	if err != nil {
		logger.Fatalf("failed to create queue: %s", err.Error())
	}
	wg := sync.WaitGroup{}
	wg.Add(50)
	cnt := int32(0)

	rand.Seed(time.Now().Unix())
	go web.Run()
	for i := 0; i < 50; i++ {
		go func() {
			time.Sleep(time.Second*time.Duration(rand.Intn(150)))
			isPublished := b.Publish("test1", "text/plain", []byte("hi"))
			if isPublished {

				atomic.AddInt32(&cnt, 1)
				logger.Trace("pub")
			}
			wg.Done()
		}()
	}
	wg.Wait()
	logger.Tracef("published: %d", cnt)
	b.Close()
}
