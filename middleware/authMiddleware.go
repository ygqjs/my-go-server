package middleware

import (
	"net/http"

	"my-go-server/database"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
  return func(ctx *gin.Context) {
    // 获取请求头中的 Token
    token := ctx.GetHeader("token")
    if token == "" {
      ctx.JSON(http.StatusUnauthorized, gin.H{
        "success": false,
        "message": "未提供 token",
      })
      ctx.Abort()
      return
    }

    // 检查 Token 是否在 token_blacklists 表中
    var count int64
    result := database.DB.Raw("SELECT COUNT(*) FROM token_blacklists WHERE token = ?", token).Scan(&count)
    if result.Error != nil {
      ctx.JSON(http.StatusInternalServerError, gin.H{
        "success": false,
        "message": "无法验证 token",
      })
      ctx.Abort()
      return
    }

    if count > 0 {
      ctx.JSON(http.StatusUnauthorized, gin.H{
        "success": false,
        "message": "token 已失效",
      })
      ctx.Abort()
      return
    }

    // 继续处理请求
    ctx.Next()
  }
}