package bus

import (
	"context"
	"github.com/annakozyreva1/person_registrator/log"
	"github.com/streadway/amqp"
)

var (
	logger = log.Logger
)

const (
	MaxConnections = 10
)

func New(url string) *Bus {
	ctx, cancel := context.WithCancel(context.Background())
	return &Bus{
		url:    url,
		tasks:  make(chan task),
		limit:  make(chan struct{}, MaxConnections),
		ctx:    ctx,
		cancel: cancel,
	}
}

type Bus struct {
	url    string
	tasks  chan task
	limit  chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

func (b *Bus) CreateQueue(name string, durable bool, autoDelete bool) error {
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

func (b *Bus) addWorkerIsPossible(task task) bool {
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

func (b *Bus) Publish(queue string, contentType string, body []byte) bool {
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

func (b *Bus) Close() {
	b.cancel()
}
