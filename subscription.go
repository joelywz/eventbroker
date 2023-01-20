package eventbroker

type Subscription struct {
	id         string
	Channel    MessageChannel
	ErrChannel chan error
}

func NewSubscription(id string) *Subscription {
	return &Subscription{
		id:         id,
		Channel:    make(chan []byte),
		ErrChannel: make(chan error),
	}
}
