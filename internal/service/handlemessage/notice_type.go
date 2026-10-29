package handlemessage

import (
	"context"
	"qqbot/internal/model"
)

func Poke(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}

func ClientStatus(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}

func FriendRecall(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}
func GroupRecall(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}

func GroupBan(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}
func GroupListBan(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}

func GroupIncrease(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}
func GroupKickMe(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}

func GroupSetAdmin(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}
func GroupUnsetAdmin(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.WebsocketMessage, event *model.Event) error {

	}
}
