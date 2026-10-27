package plug

import (
	"qqbot/internal/plug/agenttoolplug"
	"qqbot/internal/plug/eventplug"

	newlog "github.com/Lqingzhe/go-common-tool/log"
)

func RegisterPlug(logger newlog.Logger) {
	eventplug.RegisterEventPlug()
	agenttoolplug.RegisterAgentToolPlug(logger)
}
