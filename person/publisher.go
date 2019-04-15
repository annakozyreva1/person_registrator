package person

import (
	"fmt"
	"github.com/annakozyreva1/person_registrator/bus"
	"time"
)

func New(bus *bus.Bus, queue string) *Registrator {
	return &Registrator{
		bus:   bus,
		queue: queue,
	}
}

type Registrator struct {
	bus   *bus.Bus
	queue string
}

func (r *Registrator) Add(firstName string, lastName string) bool {
	message := createMessage(firstName, lastName)
	return r.bus.Publish(r.queue, "text/json", []byte(message))
}

func createMessage(firstName string, lastName string) string {
	return fmt.Sprintf("{\"first_name\": \"%s\", \"last_name\": \"%s\", \"ts\": %d}", firstName, lastName, time.Now().Unix())
}
