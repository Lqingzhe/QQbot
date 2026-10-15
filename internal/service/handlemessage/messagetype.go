package handlemessage

import (
	"context"
	"demo/internal/model"
)

//func(context.Context, chan model.Message,*model.Event) error

func Private(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}
func Group(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}
