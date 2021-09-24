package bus

import (
	"github.com/streadway/amqp"
)

func InitBus(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}
