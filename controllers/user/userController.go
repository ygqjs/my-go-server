package user

import (
	"fmt"
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
	fmt.Println("result.Error", result.Error)
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
