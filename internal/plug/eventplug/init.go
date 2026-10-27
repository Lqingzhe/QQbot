package eventplug

import (
	"qqbot/internal/model"
	"qqbot/internal/service/handlemessage"
)

func RegisterEventPlug(Config *model.Config) *handlemessage.OperateHub {
	OperatorHub := handlemessage.NewEventOperateHub(Config) // 官方插件

	return OperatorHub
}
