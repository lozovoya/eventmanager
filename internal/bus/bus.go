package bus

import (
	"EventManager/internal/model"
	"context"
)

type Call interface {
	CallToBus(ctx context.Context, call *model.Call) error
}
