package timelinelesstool

import (
	"qqbot/internal/agent/chatllm/agenttool"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

func RegisterTimelineLessToolPlug(logger newlog.Logger) {
	agenttool.RegisterOfficialTimelinessTool(logger) // 官方插件
}
