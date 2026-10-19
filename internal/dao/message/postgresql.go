package message

import (
	"context"
	"qqbot/internal/model"

	"github.com/Lqingzhe/go-common-tool/commondao"
)

func (i *MessageInfo) GetDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (gotInfo *MessageInfo, exist bool, err error) {

}
func (i *MessageInfo) SetDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (err error) {

}
func (i *MessageInfo) DeleteDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (err error) {

}
func (i *MessageInfo) UpdateDB(ctx context.Context, dbContext commondao.DBContextInterface[model.DBContext]) (exist bool, err error) {

}
