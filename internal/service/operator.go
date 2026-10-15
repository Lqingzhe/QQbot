package service

import (
	"context"
	"demo/internal/model"
	"demo/internal/service/handlemessage"
	"demo/internal/tool"
	newerror "demo/pkg/error"
	"errors"
	"log"
	"sync"
)

type Operate func(context.Context, chan model.Message, *model.Event) error

type OperateHub struct {
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
type operator struct {
	operateHub   *OperateHub
	immutableKey map[string]struct{}
	sendCh       chan model.Message
	groupList    *GroupList
	privateList  *PrivateList
	selfQQ       int64
}

var Operator *operator

func init() {
	Operator = &operator{
		operateHub: &OperateHub{
			hub: make(map[string]Operate),
		},
		groupList: &GroupList{
			List: make(map[int64]struct{}),
		},
		privateList: &PrivateList{
			List: make(map[int64]struct{}),
		},
		immutableKey: make(map[string]struct{}),
	}
}
func SetOperatorHub(sendCh chan model.Message, selfQQID int64) *operator {
	Operator.sendCh = sendCh
	Operator.selfQQ = selfQQID
	return Operator
}
func GetOperator() *operator {
	return Operator
}
func IsInGroup(groupID int64) bool {
	Operator.groupList.rwMu.RLock()
	defer Operator.groupList.rwMu.RUnlock()
	_, ok := Operator.groupList.List[groupID]
	return ok
}
func IsInPrivate(qqID int64) bool {
	Operator.privateList.rwMu.RLock()
	defer Operator.privateList.rwMu.RUnlock()
	_, ok := Operator.privateList.List[qqID]
	return ok
}
func GetSelfQQ() int64 {
	return Operator.selfQQ
}
func Register(key string, op Operate) {
	_, ok := Operator.immutableKey[key]
	if ok {
		panic("该分支不支持修改：" + key)
	}
	Operator.operateHub.rwMu.Lock()
	_, ok = Operator.operateHub.hub[key]
	if ok {
		log.Println("已覆盖原实现：" + key)
	}
	Operator.operateHub.hub[key] = op
	Operator.operateHub.rwMu.Unlock()
} //注册事件类插件

func registerBasicOp(key string, op Operate) {
	Operator.operateHub.rwMu.Lock()
	defer Operator.operateHub.rwMu.Unlock()
	_, ok := Operator.immutableKey[key]
	if ok {
		panic("重复设置：" + key)
	}
	Operator.immutableKey[key] = struct{}{}
	Operator.operateHub.hub[key] = op
}
func (o *operator) do(ctx context.Context, key string, Event *model.Event) error {
	o.operateHub.rwMu.RLock()
	op, ok := o.operateHub.hub[key]
	o.operateHub.rwMu.RUnlock()
	if !ok {
		return newerror.MakeError(newerror.CodeMethodNotAllowed, errors.New("该分支未定义行为"), newerror.LevelInfo)
	}
	return op(ctx, o.sendCh, Event)
}
func init() {
	officialPlug := handlemessage.NewOfficialPlug(IsInPrivate, IsInGroup, GetSelfQQ)
	{
		//message
		registerBasicOp(tool.GetRoutKet(model.MessagePos.String(), model.PrivateMsg.String()), officialPlug.Private)
		registerBasicOp(tool.GetRoutKet(model.MessagePos.String(), model.GroupMsg.String()), officialPlug.Group)
	}

	{
		//notice
		registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_Poke.String()), officialPlug.Poke)

		registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.ClientStatusNt.String()), officialPlug.ClientStatus)

		registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.FriendRecallNt.String()), officialPlug.FriendRecall)
		registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupRecallNt.String()), officialPlug.GroupRecall)
		{
			//GroupBan
			registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_Ban.String()), officialPlug.GroupBan)
			registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_ListBan.String()), officialPlug.GroupListBan)
		}
		registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupIncreaseNt.String()), officialPlug.GroupIncrease)
		{
			//GroupDecrease
			registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_KickMe.String()), officialPlug.GroupKickMe)
		}
		{
			//GroupAdmin
			registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Set.String()), officialPlug.GroupSetAdmin)
			registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Unset.String()), officialPlug.GroupUnsetAdmin)
		}
	}
	{
		//Request
		{
			//GroupRequest
			registerBasicOp(tool.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Add.String()), officialPlug.GroupApplyReq)
			registerBasicOp(tool.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Invite.String()), officialPlug.GroupInviteReq)
		}
	}
}
