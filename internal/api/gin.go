package api

import (
	"qqbot/internal/middleware"
	newlog "qqbot/pkg/log"

	"github.com/gin-gonic/gin"
)

type Config struct {
	logger newlog.Logger
}

func NewConfig(logger newlog.Logger) *Config {
	return &Config{
		logger: logger,
	}
}
func (c *Config) Run() {

	g := gin.New()
	g.Use(gin.Recovery(), middleware.Cors(), middleware.Log())

}
