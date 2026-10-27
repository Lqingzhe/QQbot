package handlemessage

import (
	"context"
	"qqbot/internal/model"
)

func GroupApplyReq(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.Message, event *model.Event) error {

	}
}
func GroupInviteReq(operator *Operator) Operate {
	return func(ctx context.Context, sendCh chan model.Message, event *model.Event) error {

	}
}
