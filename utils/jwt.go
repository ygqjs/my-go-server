package utils

import (
	"errors"
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

// ParseToken 解析JWT token，返回username和过期时间
// Claims 自定义声明结构体
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
func ParseToken(tokenString string) (username string, err error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	// 处理解析错误
	if err != nil {
		return "", err
	}

	// 验证token是否有效
	if !token.Valid {
		return "", errors.New("invalid token")
	}

	// 提取claims
	if claims, ok := token.Claims.(*Claims); ok {
		return claims.Username, nil
	}

	return "", errors.New("invalid claims")
}