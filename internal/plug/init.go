package plug

import (
	"qqbot/internal/agent/chatllm/agenttool"
	"qqbot/internal/model"
	"qqbot/internal/plug/agenttoolplug"
	"qqbot/internal/plug/eventplug"
	"qqbot/internal/plug/scheduletaskplug"
	"qqbot/internal/service"
	"qqbot/internal/service/handlemessage"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

func RegisterPlug(config *model.Config, logger newlog.Logger) (*handlemessage.OperateHub, *service.TaskHub, *agenttool.Tool, error) {
	eventOperator := eventplug.RegisterEventPlug(config)
	scheduleTaskHub := scheduletaskplug.RegisterScheduleTask()
	tool, err := agenttoolplug.RegisterAgentToolPlug(logger)
	if err != nil {
		return nil, nil, nil, err
	}
	return eventOperator, scheduleTaskHub, tool, nil
}
