package bus

import (
	"github.com/annakozyreva1/person_registrator/log"
	"github.com/streadway/amqp"
	"context"
)

var (
	logger = log.Logger
)

const (
	MaxConnections = 10
)

func New(url string) *bus {
	ctx, cancel := context.WithCancel(context.Background())
	return &bus{
		url: url,
		tasks: make(chan task),
		limit: make(chan struct{}, MaxConnections),
		ctx: ctx,
		cancel: cancel,
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
			go worker(b.ctx, b.url, b.tasks, b.limit)
			b.tasks <- task
		}
	default:
		{
			return false
		}
	}
	return true
}

func (b *bus) Publish(queue string, contentType string, body []byte) bool {
	task := makeTask(queue, contentType, body)
	select {
	case b.tasks <- task:
		{
		}
	default:
		{
			if !b.addWorkerIsPossible(task) {
				logger.Warning("exceed bus connections limit")
				return false
			}
		}
	}
	return task.IsSuccess()
}

func (b *bus) Close() {
	b.cancel()
}
