package handlemessage

import (
	"context"
	"qqbot/internal/model"
)

//func(context.Context, chan model.Message,*model.Event) error

func Private(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {
		if !operator.operateHub.IsInPrivate(event.UserID) {
			return nil
		}

	}
}
func Group(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {
		if !operator.operateHub.IsInGroup(event.GroupID) {
			return nil
		}

	}
}
