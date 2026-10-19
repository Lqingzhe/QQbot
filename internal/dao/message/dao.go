package message

import (
	"context"
	"qqbot/internal/model"
)

type MessageInfo struct {
	model.MessageData
	GotInfo []*model.MessageData
}

func NewMessageInfo() *MessageInfo {
	return &MessageInfo{}
}

func (i *MessageInfo) Get(ctx context.Context, dbContext *model.DBContext) (info *MessageInfo, exist bool, err error) {

}
func (i *MessageInfo) Set(ctx context.Context, dbContext *model.DBContext) (err error) {

}
func (i *MessageInfo) Delete(ctx context.Context, dbContext *model.DBContext) (err error) {

}
