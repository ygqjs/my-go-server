package models

type User struct {
	UserName string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

// 表示配置操作数据库的表名称
func (User) TableName() string {
	return "user"
}