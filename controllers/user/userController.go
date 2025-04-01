package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"my-go-server/database"
	"my-go-server/models"
	"my-go-server/utils"
)

type UserController struct {
}
func (userController UserController) Login(ctx *gin.Context) {
	// 获取表单参数username和password
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 查询数据库用户名和密码是否对应，不对应则登陆失败
	var dbUser models.User
	result := database.DB.Raw("SELECT * FROM user WHERE username = ? AND password = ? LIMIT 1", user.UserName, user.Password).Scan(&dbUser)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户名或密码错误",
		})
		return
	}

	// 生成 Token
	token, err := utils.GenerateToken(dbUser.UserName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "生成 Token 失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"token": token,
		},
	})
}

func (userController UserController) GetUserInfo(ctx *gin.Context) {
	// 获取请求头Token
	token := ctx.GetHeader("token")
	// 解析出token中的exp过期时间，如果过期了则返回信息：“token已过期”
	username, err := utils.ParseToken(token)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "token已过期",
		})
		return
	}
	// 查询数据库用户名中与username所对应的用户的用户信息
	var dbUser models.User
	result := database.DB.Raw("SELECT * FROM user WHERE username = ? LIMIT 1", username).Scan(&dbUser)
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取用户信息成功",
		"data": gin.H{
			"username": dbUser.UserName,
			"password": dbUser.Password,
		},
	})
}
