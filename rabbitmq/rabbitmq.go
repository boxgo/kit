package rabbitmq

import (
	"time"

	"github.com/streadway/amqp"
)

type (
	// RabbitMQ rabbitmq connection
	RabbitMQ struct {
		URI        string        `config:"uri" desc:"Connection uri"`
		Vhost      string        `config:"vhost" desc:"Vhost specifies the namespace of permissions, exchanges, queues and bindings on the server. Dial sets this to the path parsed from the URL."`
		ChannelMax int           `config:"channelMax" desc:"0 max channels means 2^16 - 1"`
		FrameSize  int           `config:"frameSize" desc:"0 max bytes means unlimited"`
		Heartbeat  time.Duration `config:"heartbeat" desc:"less than 1s uses the server's interval"`

		name string
		*amqp.Connection
	}
)

var (
	// Default rabbitmq
	Default = New("rabbitmq")
)

// Name config name
func (mq *RabbitMQ) Name() string {
	return mq.name
}

// ConfigWillLoad before hook
func (mq *RabbitMQ) ConfigWillLoad(context.Context) {

}

// ConfigDidLoad after hook
func (mq *RabbitMQ) ConfigDidLoad(context.Context) {
	if mq.URI == "" {
		panic("rabbitmq config is invalid")
	}

	conn, err := amqp.DialConfig(mq.URI, amqp.Config{
		Vhost:      mq.Vhost,
		ChannelMax: mq.ChannelMax,
		FrameSize:  mq.FrameSize,
		Heartbeat:  mq.Heartbeat,
	})

	if err != nil {
		panic(err)
	}

	mq.Connection = conn
}

// New options
func New(name string) *RabbitMQ {
	return &RabbitMQ{
		name: name,
	}
}
