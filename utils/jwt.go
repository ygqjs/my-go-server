package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

/**
 * 生成 JWT Token
 * @Author: 1999single
 * @Description:
 * @File: jwt
 * @Version: 1.0.0
 * @Date: 2022/5/14 20:07
 */

var jwtSecret = []byte("123456") // 替换为你的密钥
func GenerateToken(username string) (string, error) {
	// 定义 Token 的声明
	claims := jwt.MapClaims{
		"username": username, // 自定义字段
		"exp":  time.Now().Add(time.Minute * 15).Unix(), // 过期时间 (24小时)
	}
	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥签名 Token
	return token.SignedString(jwtSecret)
}