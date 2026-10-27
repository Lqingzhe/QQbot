package agenttoolplug

import (
	"qqbot/internal/agent/chatllm/agenttool"
	"qqbot/internal/plug/agenttoolplug/retunetool"
	"qqbot/internal/plug/agenttoolplug/timelinelesstool"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

func RegisterAgentToolPlug(logger newlog.Logger) (*agenttool.Tool, error) {
	tool := agenttool.NewTool(logger)
	if err := retunetool.RegisterRetuneToolPlug(tool); err != nil {
		return nil, err
	}
	if err := timelinelesstool.RegisterTimelineLessToolPlug(tool); err != nil {
		return nil, err
	}
	return tool, nil
}
