package views

import (
	"EventManager/internal/model"
	"context"
)

func SnapShot (ctx context.Context, calls []*model.Call) (*SnapshotDTO, error) {

	var queueMap = make(map[string][]*model.Call)
	for _, call := range calls {
		queueMap[call.Queue_ID] = append(queueMap[call.Queue_ID], call)
	}
	var result SnapshotDTO
	for key, calls := range queueMap {
		var queue Queue
		queue.ID = key
		queue.Calls = calls
		result.Queues = append(result.Queues, queue)
	}
	return &result, nil
}