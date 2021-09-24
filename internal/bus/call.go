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
	busConn *amqp.Connection
}

func NewCallBus(busConn *amqp.Connection) Call {
	return &rabbitBus{busConn: busConn}
}

func (r rabbitBus) CallToBus(ctx context.Context, call *model.Call) error {

	amqpChannel, err := r.busConn.Channel()
	if err != nil {
		return fmt.Errorf("Execute: %w", err)
	}
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare(
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

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
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
