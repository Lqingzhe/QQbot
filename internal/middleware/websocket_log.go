package middleware

import (
	newlog "github.com/Lqingzhe/go-common-tool/log"
)

func WebBeginLog(rawLogger newlog.Logger, traceID string) newlog.Logger {
	return rawLogger.AddTraceID(traceID)
}
func WebEndLog(logger newlog.Logger) {

}
