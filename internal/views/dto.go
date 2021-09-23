package views

import "EventManager/internal/model"

type SnapshotDTO struct {
	Queues []Queue `json:"queues"`
}

type Queue struct {
	ID string `json:"id"`
	Calls []*model.Call `json:"calls"`
}