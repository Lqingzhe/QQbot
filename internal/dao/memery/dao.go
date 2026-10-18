package memery

import (
	"context"
	"demo/internal/model"
)

type MemeryInfo struct {
	model.MemeryModel
	GotInfo []model.MemeryModel
}

func NewMemeryInfo() *MemeryInfo {
	return &MemeryInfo{}
}

func (i *MemeryInfo)Get(ctx context.Context, dbContext *model.DBContext) (info *MemeryInfo, exist bool,err error) {

}
func (i *MemeryInfo)Set(ctx context.Context,dbContext *model.DBContext) (err error) {

}
func (i *MemeryInfo)Delete(ctx context.Context, dbContext *model.DBContext) (err error) {

}
func (i *MemeryInfo)Update(ctx context.Context, dbContext *model.DBContext) (exist bool,err error) {

}