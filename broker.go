package eventbroker

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrIdExist = errors.New("id already exist")
)

type MessageChannel = chan []byte

type Broker struct {
	unregisterChannel chan Subscription
	registerChannel   chan Subscription
	broadcastChannel  MessageChannel

	subscriptions map[string]MessageChannel
}

type Message struct {
	ID    string
	Event string
	Data  interface{}
}

func New() *Broker {
	return &Broker{
		unregisterChannel: make(chan Subscription),
		registerChannel:   make(chan Subscription),
		broadcastChannel:  make(MessageChannel),
		subscriptions:     make(map[string]MessageChannel),
	}
}

func (b *Broker) Listen() {
	for {
		select {
		case s := <-b.registerChannel:
			if _, ok := b.subscriptions[s.id]; ok {
				s.ErrChannel <- ErrIdExist
			} else {
				b.subscriptions[s.id] = s.Channel
			}
		case s := <-b.unregisterChannel:
			delete(b.subscriptions, s.id)
		case msg := <-b.broadcastChannel:
			for _, v := range b.subscriptions {
				v <- msg
			}

		}
	}

}

func (b *Broker) Broadcast(msg *Message) {

	d, _ := json.Marshal(msg.Data)

	data := fmt.Sprintf("id: %s\nevent: %s\ndata: %s\n\n", msg.ID, msg.Event, string(d))

	b.broadcastChannel <- []byte(data)
}

func (b *Broker) Subscribe(id string) *Subscription {

	s := Subscription{
		id:         id,
		Channel:    make(chan []byte),
		ErrChannel: make(chan error),
	}

	b.registerChannel <- s
	return &s
}

func (b *Broker) Unsubscribe(s *Subscription) {

	b.unregisterChannel <- *s
	close(s.Channel)
	close(s.ErrChannel)
}
