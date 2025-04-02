package middleware

import (
	"github.com/gin-gonic/gin"
)

// NoCache 中间件，用于禁用缓存
func NoCache() gin.HandlerFunc {
  return func(ctx *gin.Context) {
    ctx.Writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
    ctx.Writer.Header().Set("Pragma", "no-cache")
    ctx.Writer.Header().Set("Expires", "0")
    ctx.Writer.Header().Set("Surrogate-Control", "no-store")
    ctx.Next()
  }
}