package handlemessage

import (
	"context"
	"errors"
	"log"
	"qqbot/internal/agent"
	"qqbot/internal/model"
	"qqbot/internal/tool"
	"sync"

	"github.com/Lqingzhe/go-common-tool/error"
)

type Operate func(context.Context, chan model.Message, *model.Event) error

type GroupList struct {
	List map[int64]struct{}
	rwMu sync.RWMutex
}
type PrivateList struct {
	List map[int64]struct{}
	rwMu sync.RWMutex
}
type OperateHub struct {
	temporaryHub map[string]Operate
	sumHub       map[string]Operate
	immutableKey map[string]struct{}
	rwMu         sync.RWMutex
	groupList    *GroupList
	privateList  *PrivateList
	config       *model.Config
}

func NewEventOperateHub(config *model.Config) *OperateHub {
	newHub := &OperateHub{
		temporaryHub: make(map[string]Operate),
		sumHub:       make(map[string]Operate),
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
			newHub.privateList.List[QQID.QQid] = struct{}{}
		} else {
			newHub.groupList.List[QQID.QQid] = struct{}{}
		}
	}
	return newHub
}
func (o *OperateHub) IsInGroup(groupID int64) bool {
	o.groupList.rwMu.RLock()
	defer o.groupList.rwMu.RUnlock()
	_, ok := o.groupList.List[groupID]
	return ok
}
func (o *OperateHub) IsInPrivate(qqID int64) bool {
	o.privateList.rwMu.RLock()
	defer o.privateList.rwMu.RUnlock()
	_, ok := o.privateList.List[qqID]
	return ok
}
func (o *OperateHub) GetSelfQQ() int64 {
	return o.config.Service.SelfQQ
}
func (o *OperateHub) RegisterOperator(key string, op Operate) error {
	o.rwMu.Lock()
	_, ok := o.temporaryHub[key]
	if ok {
		log.Println("operator已覆盖原实现：" + key)
	}
	o.temporaryHub[key] = op
	o.rwMu.Unlock()
	return nil
} //注册事件类插件

func (o *OperateHub) registerBasicOp(key string, op Operate) {
	o.rwMu.Lock()
	defer o.rwMu.Unlock()
	_, ok := o.immutableKey[key]
	if ok {
		panic("重复设置：" + key)
	}
	o.immutableKey[key] = struct{}{}
	o.sumHub[key] = op
}

type Operator struct {
	operateHub *OperateHub
	dbContext  *model.DBContext
	agent      *agent.Agent
}

func NewEventOperator(OperateHub *OperateHub, dbContext *model.DBContext, Agent *agent.Agent) (*Operator, error) {
	newOperator := &Operator{
		operateHub: OperateHub,
		dbContext:  dbContext,
		agent:      Agent,
	}
	newOperator.registerOfficialOP()
	err := newOperator.registerOP()
	if err != nil {
		return nil, err
	}
	return newOperator, nil
}
func (o *Operator) Do(ctx context.Context, key string, Event *model.Event, sendCh chan model.Message, errCh chan error) {
	o.operateHub.rwMu.RLock()
	op, ok := o.operateHub.sumHub[key]
	o.operateHub.rwMu.RUnlock()
	if !ok {
		errCh <- newerror.MakeError(newerror.CodeMethodNotAllowed, errors.New("该分支未定义行为: "+key), newerror.LevelInfo)
		return
	}
	go func() {
		err := op(ctx, sendCh, Event)
		if err != nil {
			errCh <- err
		}
	}()
}
func (o *Operator) registerOP() error {
	for key, op := range o.operateHub.temporaryHub {
		if _, ok := o.operateHub.immutableKey[key]; ok {
			return newerror.MakeError(newerror.CodeResourceDuplicate, errors.New("Try To Cover Official Event Plug"), newerror.LevelError)
		}
		o.operateHub.sumHub[key] = op
	}
	return nil
}
func (o *Operator) registerOfficialOP() {
	{
		//message
		o.operateHub.registerBasicOp(tool.GetRoutKet(model.MessagePos.String(), model.PrivateMsg.String()), Private(o))
		o.operateHub.registerBasicOp(tool.GetRoutKet(model.MessagePos.String(), model.GroupMsg.String()), Group(o))
	}

	{
		//notice
		o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.NotifyNt.String(), model.NotifyNt_Poke.String()), Poke(o))

		o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.ClientStatusNt.String()), ClientStatus(o))

		o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.FriendRecallNt.String()), FriendRecall(o))
		o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupRecallNt.String()), GroupRecall(o))
		{
			//GroupBan
			o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_Ban.String()), GroupBan(o))
			o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupBanNt.String(), model.GroupBan_ListBan.String()), GroupListBan(o))
		}
		o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupIncreaseNt.String()), GroupIncrease(o))
		{
			//GroupDecrease
			o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupDecreaseNt.String(), model.GroupDecrease_KickMe.String()), GroupKickMe(o))
		}
		{
			//GroupAdmin
			o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Set.String()), GroupSetAdmin(o))
			o.operateHub.registerBasicOp(tool.GetRoutKet(model.NoticePos.String(), model.GroupAdminNt.String(), model.GroupAdmin_Unset.String()), GroupUnsetAdmin(o))
		}
	}
	{
		//Request
		{
			//GroupRequest
			o.operateHub.registerBasicOp(tool.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Add.String()), GroupApplyReq(o))
			o.operateHub.registerBasicOp(tool.GetRoutKet(model.RequestPos.String(), model.GroupReq.String(), model.GroupReq_Invite.String()), GroupInviteReq(o))
		}
	}
}
