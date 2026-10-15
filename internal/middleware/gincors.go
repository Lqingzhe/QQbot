package middleware

import "github.com/gin-gonic/gin"

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的来源（开发环境可以写 *，生产环境写具体域名）
		c.Header("Access-Control-Allow-Origin", "*")

		// 允许的请求方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// 允许的请求头（关键：加上 Authorization）
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Device-ID")

		// 允许携带凭证
		c.Header("Access-Control-Allow-Credentials", "true")

		c.Next()
	}
}
