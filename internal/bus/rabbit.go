package bus

import (
	"github.com/streadway/amqp"
)

type rabbitBus struct {
	busConn *amqp.Connection
}

func InitBus(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}
