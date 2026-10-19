package memery

import (
	"context"
	"qqbot/internal/model"
)

type MemeryInfo struct {
	model.MemeryModel
	IsPrivate bool
	GotInfo []model.MemeryModel
}

func NewMemeryInfo() *MemeryInfo {
	return &MemeryInfo{}
}
func (i *MemeryInfo)SetContent(content string) *MemeryInfo {
	i.Content = content
	return i
}
func(i *MemeryInfo)SetVector(vector []float32) *MemeryInfo {
	i.Vector = vector
	return i
}

缺少置信度部分的构建函数

func (i *MemeryInfo)GetWithUser(UserID int64)*MemeryInfo {
	i.MemeryModel.UserID = UserID
	i.IsPrivate=true
	return i
}
func (i *MemeryInfo)GetWithGroup(GroupID int64)*MemeryInfo{
	i.MemeryModel.SessionID=GroupID
	i.IsPrivate=false
	return i
}

func (i *MemeryInfo) Get(ctx context.Context, dbContext *model.DBContext) (info *MemeryInfo, exist bool, err error) {

}
func (i *MemeryInfo) Set(ctx context.Context, dbContext *model.DBContext) (err error) {

}
func (i *MemeryInfo) Delete(ctx context.Context, dbContext *model.DBContext) (err error) {

}
func (i *MemeryInfo) Update(ctx context.Context, dbContext *model.DBContext) (exist bool, err error) {
}
