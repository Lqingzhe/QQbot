package agenttoolplug

import (
	"qqbot/internal/plug/agenttoolplug/retunetool"
	"qqbot/internal/plug/agenttoolplug/timelinelesstool"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

func RegisterAgentToolPlug(logger newlog.Logger) {
	retunetool.RegisterRetuneToolPlug()
	timelinelesstool.RegisterTimelineLessToolPlug(logger)
}
