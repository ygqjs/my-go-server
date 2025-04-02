package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"my-go-server/database"
	"my-go-server/models"
	"my-go-server/utils"
)

type UserController struct {
}

/**
* 登录
*/
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

/**
* 获取用户信息
*/
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

/**
* 新增用户
*/
func (userController UserController) AddUser(ctx *gin.Context) {
  // 获取新增用户的参数
  var user models.User
  if err := ctx.ShouldBindJSON(&user); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "success": false,
      "message": "参数错误",
    })
    return
  }

  // 检查用户是否已存在
  var existingUser models.User
  result := database.DB.Raw("SELECT * FROM user WHERE username = ? LIMIT 1", user.UserName).Scan(&existingUser)
  if result.RowsAffected > 0 {
    ctx.JSON(http.StatusConflict, gin.H{
      "success": false,
      "message": "用户已存在",
    })
    return
  }

  // 新增用户
  execResult := database.DB.Exec("INSERT INTO user (username, password) VALUES (?, ?)", user.UserName, user.Password)
  if execResult.Error != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "success": false,
      "message": execResult.Error.Error(),
    })
    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "success": true,
    "message": "用户新增成功",
  })
}

/**
* 删除用户
*/
func (userController UserController) DeleteUser(ctx *gin.Context) {
  // 获取删除用户的参数
  username := ctx.Query("username")

  // 检查用户是否存在
  var user models.User
  result := database.DB.Raw("SELECT * FROM user WHERE username = ? LIMIT 1", username).Scan(&user)
  if result.RowsAffected == 0 {
    ctx.JSON(http.StatusNotFound, gin.H{
      "success": false,
      "message": "用户不存在",
    })
    return
  }

  // 检查是否试图删除管理员
  if username == "admin" {
    ctx.JSON(http.StatusForbidden, gin.H{
      "success": false,
      "message": "管理员用户不能被删除",
    })
    return
  }

  // 删除用户
  execResult := database.DB.Exec("DELETE FROM user WHERE username = ?", username)
  if execResult.Error != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "success": false,
      "message": execResult.Error.Error(),
    })
    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "success": true,
    "message": "用户删除成功",
  })
}

/**
* 修改用户
*/
func (userController UserController) UpdateUser(ctx *gin.Context) {
  // 获取修改用户的参数
  var params struct {
    Username    string `json:"username"`
    NewPassword string `json:"newPassword"`
  }
  if err := ctx.ShouldBindJSON(&params); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "success": false,
      "message": "参数错误",
    })
    return
  }

  // 检查用户是否存在
  var user models.User
  result := database.DB.Raw("SELECT * FROM user WHERE username = ? LIMIT 1", params.Username).Scan(&user)
  if result.RowsAffected == 0 {
    ctx.JSON(http.StatusNotFound, gin.H{
      "success": false,
      "message": "用户不存在",
    })
    return
  }

  // 修改用户密码
  execResult := database.DB.Exec("UPDATE user SET password = ? WHERE username = ?", params.NewPassword, params.Username)
  if execResult.Error != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "success": false,
      "message": execResult.Error.Error(),
    })
    return
  }

  ctx.JSON(http.StatusOK, gin.H{
    "success": true,
    "message": "用户密码修改成功",
  })
}

/**
* 查询用户列表
*/
func (userController UserController) UserList(ctx *gin.Context) {
  // 获取查询参数
  current := ctx.Query("current")       // 当前分页
  pageSize := ctx.Query("pageSize")    // 分页大小
  id := ctx.Query("id")                // 用户 ID（模糊查询）
  username := ctx.Query("username")    // 用户名（模糊查询）
  sex := ctx.Query("sex")              // 性别
  address := ctx.Query("address")      // 地址（模糊查询）

  // 校验分页参数是否为大于 1 的整数
  if current == "" || current == "0" {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "success": false,
      "message": "分页参数 current 必须是大于 1 的整数",
    })
    return
  }
  if pageSize == "" || pageSize == "0" {
    ctx.JSON(http.StatusBadRequest, gin.H{
      "success": false,
      "message": "分页参数 pageSize 必须是大于 1 的整数",
    })
    return
  }

  // 转换分页参数为整数
  currentPage, _ := strconv.Atoi(current)
  size, _ := strconv.Atoi(pageSize)

  // 构建查询条件
  query := "SELECT * FROM user WHERE 1=1"
  countQuery := "SELECT COUNT(*) FROM user WHERE 1=1"
  args := []interface{}{}

  if id != "" {
    query += " AND id LIKE ?"
    countQuery += " AND id LIKE ?"
    args = append(args, "%"+id+"%")
  }
  if username != "" {
    query += " AND username LIKE ?"
    countQuery += " AND username LIKE ?"
    args = append(args, "%"+username+"%")
  }
  if sex != "" {
    query += " AND sex = ?"
    countQuery += " AND sex = ?"
    args = append(args, sex)
  }
  if address != "" {
    query += " AND address LIKE ?"
    countQuery += " AND address LIKE ?"
    args = append(args, "%"+address+"%")
  }

  // 查询符合条件的总条数
  var total int64
  countResult := database.DB.Raw(countQuery, args...).Scan(&total)
  if countResult.Error != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "success": false,
      "message": countResult.Error.Error(),
    })
    return
  }

  // 查询分页数据
  query += " LIMIT ? OFFSET ?"
  args = append(args, size, (currentPage-1)*size)

  var users []models.User
  result := database.DB.Raw(query, args...).Scan(&users)
  if result.Error != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
      "success": false,
      "message": result.Error.Error(),
    })
    return
  }

  // 返回结果
  ctx.JSON(http.StatusOK, gin.H{
    "success": true,
    "message": "查询用户列表成功",
    "data": gin.H{
      "data":    users,
      "current": currentPage,
      "total":   total,
    },
  })
}