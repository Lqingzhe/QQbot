package memery

import (
	"context"
	"demo/internal/model"

	"github.com/Lqingzhe/go-common-tool/commondao"
)

func (i *MemeryInfo) GetDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (info *MemeryInfo, exist bool, err error) {

}
func (i *MemeryInfo) SetDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (err error) {

}
func (i *MemeryInfo) DeleteDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (err error) {

}
func (i *MemeryInfo) UpdateDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (exist bool, err error) {

}
