package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// 数据库连接基本信息
		var (
			username  = "root"
			password  = "12345678"
			ipAddress = "127.0.0.1"
			port      = 3306
			dbName    = "go_test"
			charset   = "utf8mb4"
		)
    // 配置数据库连接信息
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", username, password, ipAddress, port, dbName, charset)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("无法连接到数据库: %v", err)
    }
    DB = db
    fmt.Println("------------------------------数据库连接成功")
}