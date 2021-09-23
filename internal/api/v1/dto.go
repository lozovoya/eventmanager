package v1

import "EventManager/internal/model"

type CallEventDTO struct {
	EventID string `json:"event_id"`
	Type string `json:"event_type"`
	Timestamp string `json:"timestamp"`
	Data model.Call `json:"data"`
}

type SnapShotDTO struct {
	Queues []*Queue `json:"queues"`
}

type Queue struct {
	ID string `json:"id"`
	Calls []*model.Call `json:"calls"`
}