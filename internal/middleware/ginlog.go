package middleware

import (
	"errors"
	"qqbot/pkg/log"
	"time"

	"uuid"

	"github.com/gin-gonic/gin"
)

func Log(rawLogger newlog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		beginTime := time.Now()
		trace := uuid.NewV4().String()
		logger := rawLogger.AddTraceID(trace).AddIp(c.ClientIP()).AddOperation(c.FullPath())
		c.Set("logger", logger)
		c.Set("trace", trace)

		c.Next()
		RawLogger, exist := c.Get("logger")
		if !exist {
			logger.AddLatencyAndTime(beginTime).
				AddLevel(newerror.LevelError).
				AddError(errors.New("Middleware Can't Get Logger"), newerror.CodeInternalError).
				Log()
			return
		}
		logger2, ok := RawLogger.(newlog.Logger)
		if !ok {
			logger.AddLatencyAndTime(beginTime).
				AddLevel(newerror.LevelError).
				AddError(errors.New("Middleware's Logger's Type Is Wrong"), newerror.CodeInternalError).
				Log()
			return
		}
		logger2.AddLatencyAndTime(beginTime)
		rawId, exist := c.Get("user_id")
		if exist {
			userId, ok := rawId.(int64)
			if ok {
				logger2 = logger2.AddUserID(userId)
			} else {
				logger2.AddLevel(newerror.LevelError).
					AddError(errors.New("UserID Can't Translate Into Int64"), newerror.CodeInternalError).
					Log()
				return

			}
		}
		logger2.Log()
	}
}
