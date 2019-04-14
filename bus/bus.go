package bus

import (
	"github.com/annakozyreva1/person_registrator/log"
	"github.com/streadway/amqp"
	"context"
)

const (
	MaxConnCount = 10
)

var (
	logger = log.Logger
)

func NewBus()*bus {
	return &bus{

	}
}

type bus struct {
	url   string
	tasks chan task
	limit chan struct{}
	ctx context.Context
	cancel context.CancelFunc
}

func (b *bus) CreateQueue(name string, durable bool, autoDelete bool) error {
	conn, err := amqp.Dial(b.url)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (b *bus) addWorkerIsPossible(task task) bool {
	select {
	case b.limit <- struct{}{}:
		{
			b.tasks <- task
			go worker(b.ctx, b.url, b.tasks, b.limit)
		}
	default:
		{
			return false
		}
	}
	return true
}

func (b *bus) Publish(queue string, contentType string, body []byte) bool {
	task := task{
		Queue:       queue,
		ContentType: contentType,
		Body:        body,
	}
	select {
	case b.tasks <- task:
		{
		}
	default:
		{
			if !b.addWorkerIsPossible(task) {
				logger.Warning("exceed bus connections limit")
			}
		}
	}
	return task.IsSuccess()
}

func (b *bus) Close() {
	b.cancel()
}
