package service

import (
	"qqbot/internal/agent"
	"qqbot/internal/model"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

//type ReqModel struct {
//	Event   *model.Event
//	TraceID string
//}
//
//func NewReqModel(event *model.Event, traceID string) ReqModel {
//	return ReqModel{
//		Event:   event,
//		TraceID: traceID,
//	}
//}

//type consumer struct {
//	ch           chan ReqModel
//	cancelWorker context.CancelFunc
//}
//
//func (c *consumer) Stop() {
//	c.cancelWorker()
//}
//
//type messageList struct {
//	List map[int64]*consumer
//	RWMu sync.RWMutex
//}
//
//func newMessageList(cap ...int) *messageList {
//	list := &messageList{}
//	if len(cap) == 0 {
//		list.List = make(map[int64]*consumer)
//	} else {
//		list.List = make(map[int64]*consumer, cap[0])
//	}
//	return list
//}
//func (m *messageList) Add(id int64, consumer *consumer) {
//	m.RWMu.Lock()
//	defer m.RWMu.Unlock()
//	m.List[id] = consumer
//}
//func (m *messageList) Remove(id int64) {
//	m.RWMu.Lock()
//	defer m.RWMu.Unlock()
//	c, ok := m.Get(id)
//	if ok {
//		c.Stop()
//		delete(m.List, id)
//	}
//}
//func (m *messageList) Get(id int64) (*consumer, bool) {
//	m.RWMu.RLock()
//	c, ok := m.List[id]
//	m.RWMu.RUnlock()
//	return c, ok
//}

type Service struct {
	logCh        chan newlog.Logger
	sendCh       chan model.Message
	config       *model.Config
	operatorHub  *operator
	scheduleTask *taskHub
	agent        *agent.Agent
}

func NewService(config *model.Config, operatorHub *operator, taskHub *taskHub, agent *agent.Agent) *Service {
	return &Service{
		logCh:        make(chan newlog.Logger, config.Service.ChannelLength),
		config:       config,
		operatorHub:  operatorHub,
		scheduleTask: taskHub,
		agent:        agent,
	}
}
func (s *Service) SetSendCh(sendCh chan model.Message) {
	s.sendCh = sendCh
}
