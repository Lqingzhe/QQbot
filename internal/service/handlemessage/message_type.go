package handlemessage

import (
	"context"
	"qqbot/internal/model"
)

//func(context.Context, chan model.Message,*model.Event) error

func Private(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.Message, event *model.Event) error {

	}
}
func Group(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.Message, event *model.Event) error {

	}
}
