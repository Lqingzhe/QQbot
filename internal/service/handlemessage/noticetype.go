package handlemessage

import (
	"context"
	"demo/internal/model"
)

func Poke(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}

func ClientStatus(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}

func FriendRecall(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}
func GroupRecall(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}

func GroupBan(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}
func GroupListBan(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}

func GroupIncrease(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}
func GroupKickMe(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}

func GroupSetAdmin(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}
func GroupUnsetAdmin(ctx context.Context, sendCh chan model.Message, Event *model.Event) error {

}