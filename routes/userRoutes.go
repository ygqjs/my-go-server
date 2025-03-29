package routes

import (
	"my-go-server/controllers/user" // 确保路径正确

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/login", user.UserController{}.Login) // 调用导出的 Login 方法
	}
}