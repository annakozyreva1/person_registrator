package bus

import (
	"context"
	"github.com/streadway/amqp"
	"time"
)

const (
	MaxPublishTries       = 3
	TryTimeout            = time.Second * 10
	IdleConnectionTimeout = time.Minute
)

func connect(url string) (*amqp.Connection, *amqp.Channel,  chan amqp.Confirmation, chan *amqp.Error, error) {
	var conn *amqp.Connection
	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, nil, nil,  err
	}
	ch.Confirm(false)
	conf := make(chan amqp.Confirmation)
	ch.NotifyPublish(conf)
	 cl := make(chan *amqp.Error)
	ch.NotifyClose(cl)
	return conn, ch, conf, cl, nil
}

func publish(ch *amqp.Channel, queue string, contentType string, body []byte) error {
	return ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  contentType,
			Body:         body,
		})
}

func worker(ctx context.Context, url string, tasks chan task, limit chan struct{}) {
	logger.Debug("started bus worker")
	pub := newPublisher(url)
	defer func() {
		pub.Close()
		limit <- struct{}{}
		logger.Debug("closed bus worker")
	}()
	var err error
	for {
		select {
		case task := <-tasks:
			{
				for try := 1; try <= MaxPublishTries; try++ {
					err = pub.Publish(task.Queue, task.ContentType, task.Body)
					if err == nil {
						break
					}
					time.Sleep(TryTimeout)
				}
				if err != nil {
					logger.Errorf("failed to publish in %s: %s", task.Queue, err.Error())
					task.Failure()
				} else {
					task.Success()
				}
			}
		case <-time.After(IdleConnectionTimeout):
			{
				break
			}
		case <-ctx.Done():
			{
				break

			}
		}
	}
}