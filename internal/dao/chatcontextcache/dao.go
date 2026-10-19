package chatcache

import (
	"context"
	"qqbot/internal/model"
)

type ContextInfo struct {
	model.ChatCacheModel
	GotInfo *model.ChatCacheModel //model
}

func NewContextInfo(sessionID int64) *ContextInfo {
	return &ContextInfo{
		SessionID: sessionID,
	}
}
func (i *ContextInfo) IsGroup() *ContextInfo {
	i.ChatCacheModel.IsPrivate = false
	return i
}
func (i *ContextInfo) IsPrivate() *ContextInfo {
	i.ChatCacheModel.IsPrivate = true
	return i
}
func (i *ContextInfo) SetData(data []byte) *ContextInfo {
	i.ChatCacheModel.Data = data
	return i
}

func (i *ContextInfo) Get(ctx context.Context, dbContext *model.DBContext) (info *ContextInfo, exist bool, err error) {

}
func (i *ContextInfo) Set(ctx context.Context, dbContext *model.DBContext) (err error) {

}

// func (i *ContextInfo) Delete(ctx context.Context, dbContext *model.DBContext) (err error) {
//
// }
func (i *ContextInfo) Update(ctx context.Context, dbContext *model.DBContext) (exist bool, err error) {

}
