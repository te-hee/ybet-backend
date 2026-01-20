package dispatcher

import (
	"encoding/json"
	"log"
	"messageService/config"

	"github.com/nats-io/nats.go/jetstream"
)

type Event struct {
	Subject string
	Data    any
}

type Dispatcher struct {
	js     jetstream.JetStream
	events chan Event
}

func NewDispatcher(js jetstream.JetStream) *Dispatcher {
	return &Dispatcher{
		js:     js,
		events: make(chan Event, *config.CustomBuffer),
	}
}

func (d *Dispatcher) Start() {
	for {
		event := <-d.events

		payload, err := json.Marshal(event.Data)
		if err != nil {
			log.Printf("failed to marshal json TwT: %v", err)
		}
		ackF, err := d.js.PublishAsync(event.Subject, payload)
		if err != nil {
			log.Printf("failed to publish on NATS :c : %v", err)
		}

		select {
		case ack := <-ackF.Ok():
			log.Printf("Published msg with sequence number %d on stream %q", ack.Sequence, ack.Stream)
		case err := <-ackF.Err():
			log.Println(err)
		}
	}
}

func (d *Dispatcher) Close() {
	close(d.events)
}

func (d *Dispatcher) Publish(subject string, data any) {
	d.events <- Event{Subject: subject, Data: data}
}
