package bus

import (
	"context"
	"github.com/streadway/amqp"
	"net"
	"time"
)

const (
	MaxPublishTries       = 3
	TryTimeout            = time.Second * 10
	IdleConnectionTimeout = time.Minute
)

func connect(url string) (*amqp.Connection, *amqp.Channel, error) {
	var conn *amqp.Connection
	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	return conn, ch, nil
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
	var conn *amqp.Connection
	var ch *amqp.Channel
	defer func() {
		if ch != nil {
			ch.Close()
		}
		if conn != nil {
			conn.Close()
		}
		limit <- struct{}{}
	}()
	var err error
	for {
		select {
		case task := <-tasks:
			{
				for try := 1; try <= MaxPublishTries; try++ {
					if ch == nil {
						conn, ch, err = connect(url)
						if err != nil {
							continue
						}
					}
					err = publish(ch, task.Queue, task.ContentType, task.Body)
					if err == nil {
						break
					}
					if _, ok := err.(*net.OpError); ok {
						conn = nil
						ch = nil
					}
					time.Sleep(TryTimeout)
				}
				if err != nil {
					logger.Errorf("failed publish in %s: %s", task.Queue, err.Error())
					task.Failure()
				} else {
					logger.Tracef("published: %+v", task.Body)
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
	logger.Debug("closed bus worker")
}