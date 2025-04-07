package routes

import (
	"my-go-server/controllers/user"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/login", user.UserController{}.Login) // 调用导出的 Login 方法
		userRoutes.GET("/logout", user.UserController{}.Logout) // 调用导出的 Login 方法
		userRoutes.GET("/user-info", user.UserController{}.GetUserInfo)
	}
	r.POST("/users", user.UserController{}.AddUser)
	r.DELETE("/users", user.UserController{}.DeleteUser)
	r.PUT("/users", user.UserController{}.UpdateUser)
	r.GET("/users", user.UserController{}.UserList)
}