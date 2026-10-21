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

type Task func(context.Context, chan model.Message)
type taskHub struct {
	taskInfo map[string]Task
	ctxHub   map[string]context.CancelFunc
	mutex    sync.Mutex
}

func (t *taskHub) Do(sendCh chan model.Message) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	for taskName, task := range t.taskInfo {
		ctx, cancel := context.WithCancel(context.Background())
		task(ctx, sendCh)
		t.ctxHub[taskName] = cancel
	}
}

var TaskWorker = &taskHub{
	taskInfo: make(map[string]Task),
	ctxHub:   make(map[string]context.CancelFunc),
}

func GetAndInitTaskHub(config *model.Config) *taskHub {
	return TaskWorker
}
func GetTaskHub() *taskHub {
	return TaskWorker
}
func RegisterScheduleTask(taskName string, task Task) error {
	TaskWorker.mutex.Lock()
	defer TaskWorker.mutex.Unlock()
	_, ok := TaskWorker.taskInfo[taskName]
	if ok {
		return errors.New("The Schedule Task Name is Exist: " + taskName)
	}
	TaskWorker.taskInfo[taskName] = task
	return nil
}
func registerBasicScheduleTask(taskName string, task Task) {
	TaskWorker.mutex.Lock()
	defer TaskWorker.mutex.Unlock()
	_, ok := TaskWorker.taskInfo[taskName]
	if ok {
		log.Println("scheduleTask已覆盖原实现：" + taskName)
	}
	TaskWorker.taskInfo[taskName] = task
}
func init() {
	registerBasicScheduleTask("pingTest", func(ctx context.Context, sendCh chan model.Message) {
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

}
