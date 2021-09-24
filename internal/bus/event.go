package bus

import (
	"EventManager/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func NewEventBus(busConn *amqp.Connection) Event {
	return &rabbitBus{busConn: busConn}
}

func (r rabbitBus) EventToBus(ctx context.Context, event *model.Event) error {

	amqpChannel, err := r.busConn.Channel()
	if err != nil {
		return fmt.Errorf("EventToBus: %w", err)
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
		return fmt.Errorf("EventToBus: %w", err)
	}
	body, err := json.Marshal(&event)
	if err != nil {
		return fmt.Errorf("EventToBus: %w", err)
	}

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return fmt.Errorf("EventToBus: %w", err)
	}

	log.Printf("message sent %s, type %s", event.EventID, event.Type)
	return nil
}
