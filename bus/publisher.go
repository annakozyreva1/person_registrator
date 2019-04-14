package bus

import (
	"github.com/streadway/amqp"
	"errors"
)

func newPublisher(url string) *publisher {
	return &publisher{
		url: url,
		closed: make(chan *amqp.Error, 1),
	}
}

type publisher struct {
	url string
	conn *amqp.Connection
	ch *amqp.Channel
	confirmed chan amqp.Confirmation
	closed chan *amqp.Error
}

func (p *publisher) connect() error {
	var err error
	p.conn, err = amqp.Dial(p.url)
	if err != nil {
		return err
	}
	p.ch, err = p.conn.Channel()
	if err != nil {
		p.conn.Close()
		p.conn = nil
		return err
	}
	p.ch.Confirm(false)
	p.confirmed = make(chan amqp.Confirmation)
	p.ch.NotifyPublish(p.confirmed)
	p.closed = make(chan *amqp.Error, 1)
	p.ch.NotifyClose(p.closed)
	return nil
}

func (p *publisher) Close() {
	if p.ch != nil {
		p.ch.Close()
		p.ch = nil
	}
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
}

func (p *publisher) isClosed() bool {
	select {
	case <- p.closed:
		{
			p.Close()
			return true
		}
	default:
		{

		}
	}
	return false
}

func (p *publisher) Publish(queue string, contentType string, body []byte) error {
	if p.ch == nil || p.isClosed() {
		if err := p.connect(); err != nil  {
			return err
		}
	}
	err := p.ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  contentType,
			Body:         body,
		})
	if err == nil {
		c := <- p.confirmed
		if c.Ack {
			return nil
		}
		return errors.New("message is not sent")
	}
	return err
}
