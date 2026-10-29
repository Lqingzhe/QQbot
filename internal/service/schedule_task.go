package service

import (
	"context"
	"errors"
	"log"
	"qqbot/internal/model"
	"sync"
	"time"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

type Task func(context.Context, chan model.WebsocketMessage)
type TaskHub struct {
	taskInfo map[string]Task
	ctxHub   map[string]context.CancelFunc
	mutex    sync.Mutex
}

func (t *TaskHub) Do(sendCh chan model.WebsocketMessage) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for taskName, task := range t.taskInfo {
		ctx, cancel := context.WithCancel(context.Background())
		go task(ctx, sendCh)
		t.ctxHub[taskName] = cancel
	}
}

func NewScheduleTaskHub() *TaskHub {
	newTaskHub := &TaskHub{
		taskInfo: make(map[string]Task),
		ctxHub:   make(map[string]context.CancelFunc),
	}
	newTaskHub.registerBasicScheduleTask("pingTest", func(ctx context.Context, sendCh chan model.WebsocketMessage) {
		ticker := time.NewTicker(24 * time.Hour)
		for {
			select {
			case <-ctx.Done():
				log.Println("已退出pingTest")
				return
			case <-ticker.C:
				newlog.SetLabel("缺少发送消息的构建")
				//sendCh<- model.Message{
				//MessageType: websocket.TextMessage,
				//}
			}
		}
	})

	return newTaskHub
}

func (t *TaskHub) RegisterScheduleTask(taskName string, task Task) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	_, ok := t.taskInfo[taskName]
	if ok {
		return errors.New("The Schedule Task Name is Exist: " + taskName)
	}
	t.taskInfo[taskName] = task
	return nil
}
func (t *TaskHub) registerBasicScheduleTask(taskName string, task Task) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	_, ok := t.taskInfo[taskName]
	if ok {
		log.Println("scheduleTask已覆盖原实现：" + taskName)
	}
	t.taskInfo[taskName] = task
}
