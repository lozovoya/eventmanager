package evSender

import (
	"EventManager/internal/model"
	"log"
)

type CallEventDTO struct {
	EventID   string     `json:"event_id"`
	Type      string     `json:"event_type"`
	Timestamp string     `json:"timestamp"`
	Data      model.Call `json:"data"`
}

func main() {
	log.Println("dfsfsd")
}
