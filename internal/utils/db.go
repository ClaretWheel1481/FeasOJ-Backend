package utils

import (
	"fmt"
	"log"
	"src/internal/config"
	"src/internal/global"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 创建管理员
func InitAdminAccount() (string, string, string, string, int) {
	var adminUsername string
	var adminPassword string
	var adminEmail string
	log.Println("[FeasOJ] Please input the administrator account configuration: ")
	fmt.Print("[FeasOJ] Username: ")
	fmt.Scanln(&adminUsername)
	fmt.Print("[FeasOJ] Password: ")
	fmt.Scanln(&adminPassword)
	fmt.Print("[FeasOJ] Email: ")
	fmt.Scanln(&adminEmail)

	return adminUsername, EncryptPassword(adminPassword), adminEmail, uuid.New().String(), 1
}

// 创建表
func InitTable() bool {
	global.DB.AutoMigrate(
		&global.User{},
		&global.Problem{},
		&global.SubmitRecord{},
		&global.Discussion{},
		&global.Comment{},
		&global.TestCase{},
		&global.Competition{},
		&global.UserCompetitions{},
	)
	return true
}

// 返回数据库连接对象
func ConnectSql() *gorm.DB {
	dsn := config.LoadSqlConfig()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("[FeasOJ] Database connection failed, please go to config.xml manually to configure.")
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("[FeasOJ] Failed to get generic database object.")
		return nil
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.MaxLifeTime * time.Second)
	return db
}

// 根据用户名获取用户信息
func SelectUser(username string) global.User {
	var user global.User
	ConnectSql().Where("username = ?", username).First(&user)
	return user
}
