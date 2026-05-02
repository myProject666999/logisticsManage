package models

import (
	"fmt"
	"log"
	"logisticsManage/config"
	"logisticsManage/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dbConfig := config.AppConfig.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)

	log.Println("数据库连接成功")
}

func AutoMigrate() {
	DB.AutoMigrate(
		&User{},
		&DeliveryMan{},
		&Admin{},
		&Site{},
		&Order{},
		&Logistics{},
		&Delivery{},
		&Reimbursement{},
		&Warehouse{},
		&WarehouseItem{},
		&Settlement{},
	)
	log.Println("数据库迁移完成")
}

func SeedData() {
	var adminCount int
	DB.Model(&Admin{}).Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, _ := utils.HashPassword("123456")
		admin := &Admin{
			Username: "admin",
			Password: hashedPassword,
			RealName: "超级管理员",
			Phone:    "13800138000",
			Role:     2,
			Status:   1,
		}
		if err := DB.Create(admin).Error; err == nil {
			log.Println("默认管理员账号创建成功: admin / 123456")
		} else {
			log.Printf("创建默认管理员失败: %v", err)
		}
	} else {
		log.Println("管理员数据已存在，跳过种子数据初始化")
	}

	var userCount int
	DB.Model(&User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, _ := utils.HashPassword("123456")
		user := &User{
			Username: "user1",
			Password: hashedPassword,
			Phone:    "13900139001",
			RealName: "测试用户",
			Status:   1,
		}
		if err := DB.Create(user).Error; err == nil {
			log.Println("默认用户账号创建成功: user1 / 123456")
		}
	}

	var deliveryCount int
	DB.Model(&DeliveryMan{}).Count(&deliveryCount)
	if deliveryCount == 0 {
		hashedPassword, _ := utils.HashPassword("123456")
		deliveryMan := &DeliveryMan{
			Username:     "delivery1",
			Password:     hashedPassword,
			Phone:        "13700137001",
			RealName:     "测试配送员",
			Status:       1,
			VerifyStatus: 1,
		}
		if err := DB.Create(deliveryMan).Error; err == nil {
			log.Println("默认配送员账号创建成功: delivery1 / 123456")
		}
	}

	log.Println("种子数据初始化完成")
}
