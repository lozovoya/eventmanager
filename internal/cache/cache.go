package cache

import (
	"EventManager/internal/model"
	"context"
)

type Call interface {
	CallToCache (ctx context.Context, call *model.Call) error
	CallFromCache (ctx context.Context, queueID, callID string) error
	GetCallsSnapshot (ctx context.Context) ([]*model.Call, error)
}
