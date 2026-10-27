package handlemessage

import (
	"context"
	"errors"
	"log"
	"qqbot/internal/model"
	"qqbot/internal/tool"
	"sync"

	"github.com/Lqingzhe/go-common-tool/error"
)

type Operate func(context.Context, chan model.Message, *model.Event) error

type operateHub struct {
	hub  map[string]Operate
	rwMu sync.RWMutex
}
type GroupList struct {
	List map[int64]struct{}
	rwMu sync.RWMutex
}
type PrivateList struct {
	List map[int64]struct{}
	rwMu sync.RWMutex
}
type Operator struct {
	operateHub   *operateHub
	immutableKey map[string]struct{}
	groupList    *GroupList
	privateList  *PrivateList
	config       *model.Config
}

func NewEventOperator(config *model.Config) *Operator {
	newOperator := &Operator{
		operateHub:   &operateHub{hub: make(map[string]Operate)},
		immutableKey: make(map[string]struct{}),
		groupList: &GroupList{
			List: make(map[int64]struct{}),
		},
		privateList: &PrivateList{
			List: make(map[int64]struct{}),
		},
		config: config,
	}
	for _, QQID := range config.ResponserList {
		if QQID.IsPrivate {
			newOperator.privateList.List[QQID.QQid] = struct{}{}
		} else {
			newOperator.groupList.List[QQID.QQid] = struct{}{}
		}
	}
	newOperator.registerOfficialOP()
	return newOperator
}
func (o *Operator) IsInGroup(groupID int64) bool {
	o.groupList.rwMu.RLock()
	defer o.groupList.rwMu.RUnlock()
	_, ok := o.groupList.List[groupID]
	return ok
}
func (o *Operator) IsInPrivate(qqID int64) bool {
	o.privateList.rwMu.RLock()
	defer o.privateList.rwMu.RUnlock()
	_, ok := o.privateList.List[qqID]
	return ok
}
func (o *Operator) GetSelfQQ() int64 {
	return o.config.Service.SelfQQ
}
func (o *Operator) RegisterOperator(key string, op Operate) error {
	_, ok := o.immutableKey[key]
	if ok {
		return errors.New("该分支不支持修改：" + key)
	}
	o.operateHub.rwMu.Lock()
	_, ok = o.operateHub.hub[key]
	if ok {
		log.Println("operator已覆盖原实现：" + key)
	}
	o.operateHub.hub[key] = op
	o.operateHub.rwMu.Unlock()
	return nil
} //注册事件类插件

func (o *Operator) registerBasicOp(key string, op Operate) {
	o.operateHub.rwMu.Lock()
	defer o.operateHub.rwMu.Unlock()
	_, ok := o.immutableKey[key]
	if ok {
		panic("重复设置：" + key)
	}
	o.immutableKey[key] = struct{}{}
	o.operateHub.hub[key] = op
}
func (o *Operator) Do(ctx context.Context, key string, Event *model.Event, sendCh chan model.Message) error {
	o.operateHub.rwMu.RLock()
	op, ok := o.operateHub.hub[key]
	o.operateHub.rwMu.RUnlock()
	if !ok {
		return newerror.MakeError(newerror.CodeMethodNotAllowed, errors.New("该分支未定义行为"), newerror.LevelInfo)
	}
	return op(ctx, sendCh, Event)
}
func (o *Operator) registerOfficialOP() {
	{
		//message
		o.registerBasicOp(tool.GetRoutKet(model.MessagePos.String(), model.PrivateMsg.String()), Private(o))
		o.registerBasicOp(tool.GetRoutKet(model.MessagePos.String(), model.GroupMsg.String()), Group(o))
	}

	{
		//notice
		o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_Poke.String()), Poke(o))

		o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.ClientStatusNt.String()), ClientStatus(o))

		o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.FriendRecallNt.String()), FriendRecall(o))
		o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupRecallNt.String()), GroupRecall(o))
		{
			//GroupBan
			o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_Ban.String()), GroupBan(o))
			o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_ListBan.String()), GroupListBan(o))
		}
		o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupIncreaseNt.String()), GroupIncrease(o))
		{
			//GroupDecrease
			o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_KickMe.String()), GroupKickMe(o))
		}
		{
			//GroupAdmin
			o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Set.String()), GroupSetAdmin(o))
			o.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Unset.String()), GroupUnsetAdmin(o))
		}
	}
	{
		//Request
		{
			//GroupRequest
			o.registerBasicOp(tool.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Add.String()), GroupApplyReq(o))
			o.registerBasicOp(tool.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Invite.String()), GroupInviteReq(o))
		}
	}
}
