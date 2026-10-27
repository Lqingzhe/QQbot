package eventplug

import (
	"qqbot/internal/model"
	"qqbot/internal/service/handlemessage"

)

func RegisterEventPlug(Config *model.Config) *handlemessage.Operator{
	Operator :=handlemessage.NewEventOperator(Config) // 官方插件

	return Operator
}
