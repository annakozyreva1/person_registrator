package bus

import (
	"context"
	"github.com/streadway/amqp"
	"net"
	"time"
	"reflect"
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
	logger.Debug("started bus worker")
	defer func() {
		if ch != nil {
			ch.Close()
		}
		if conn != nil {
			conn.Close()
		}
		limit <- struct{}{}
		logger.Debug("closed bus worker")
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
							logger.Debug("err try", try)
							continue
						}
					}
					err = publish(ch, task.Queue, task.ContentType, task.Body)
					if err == nil {
						break
					}
					conn = nil
					ch = nil
					logger.Debugf("%+v", reflect.TypeOf(err))
					if _, ok := err.(*net.OpError); ok {
						conn = nil
						ch = nil
						logger.Debug("err try", try)
					}
					time.Sleep(TryTimeout)
				}
				if err != nil {
					logger.Errorf("failed to publish in %s: %s", task.Queue, err.Error())
					task.Failure()
				} else {
					task.Success()
					logger.Debug("pub")
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