package scheduletaskplug

import "qqbot/internal/service"

func RegisterScheduleTask() *service.TaskHub {
	scheduleTaskHub := service.NewScheduleTaskHub()

	return scheduleTaskHub
}
