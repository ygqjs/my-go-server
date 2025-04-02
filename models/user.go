package models

type User struct {
	Id string `gorm:"column:id" json:"id"`
	UserName string `gorm:"column:username" json:"username"`
	Address string `gorm:"column:address" json:"address"`
	Sex string `gorm:"column:sex" json:"sex"`
	Password string `gorm:"column:password" json:"password"`
}
