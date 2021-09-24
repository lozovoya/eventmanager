package bus

import (
	"EventManager/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type rabbitBus struct {
	amqpChannel *amqp.Channel
}

func InitBus(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func NewBus(amqpChannel *amqp.Channel) Call {
	return &rabbitBus{amqpChannel: amqpChannel}
}

func (r rabbitBus) CallToBus(ctx context.Context, call *model.Call) error {
	queue, err := r.amqpChannel.QueueDeclare(
		"calls",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return fmt.Errorf("CallToBus: %w", err)
	}
	body, err := json.Marshal(&call)
	if err != nil {
		return fmt.Errorf("CallToBus: %w", err)
	}

	err = r.amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return fmt.Errorf("CallToBus: %w", err)
	}

	log.Printf("message sent %s", call.CallID)
	return nil
}
