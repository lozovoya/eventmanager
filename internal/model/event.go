package model

type Event struct {
	EventID string `json:"event_id"`
	Type string `json:"event_type"`
	Timestamp string `json:"timestamp"`
	Data interface{} `json:"data"`
}
