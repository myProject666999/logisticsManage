package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 用户表
type User struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	Phone    string `gorm:"unique_index"`
	Email    string
	RealName string
	Address  string
	Status   int `gorm:"default:1"` // 1:正常 0:禁用
}

// 配送员表
type DeliveryMan struct {
	gorm.Model
	Username     string `gorm:"unique_index;not null"`
	Password     string `gorm:"not null"`
	Phone        string `gorm:"unique_index"`
	RealName     string
	IDCard       string
	VehicleNo    string // 车牌号
	VehicleType  int    // 车辆类型：1电动车 2摩托车 3货车
	SiteID       uint
	Status       int `gorm:"default:1"` // 1:正常 0:禁用
	VerifyStatus int `gorm:"default:0"` // 0:待审核 1:已通过 2:已拒绝
}

// 管理员表
type Admin struct {
	gorm.Model
	Username string `gorm:"unique_index;not null"`
	Password string `gorm:"not null"`
	RealName string
	Phone    string
	Role     int `gorm:"default:1"` // 1:普通管理员 2:超级管理员
	Status   int `gorm:"default:1"` // 1:正常 0:禁用
}

// 站点表
type Site struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Province  string
	City      string
	District  string
	Address   string `gorm:"not null"`
	Phone     string
	Manager   string
	Status    int `gorm:"default:1"` // 1:正常 0:关闭
}

// 订单表
type Order struct {
	gorm.Model
	OrderNo       string `gorm:"unique_index;not null"`
	UserID        uint
	SenderName    string
	SenderPhone   string
	SenderAddress string
	ReceiverName  string
	ReceiverPhone string
	ReceiverAddr  string
	GoodsType     string
	GoodsWeight   float64
	GoodsAmount   int
	InsuranceFee  float64
	TotalFee      float64
	Status        int `gorm:"default:0"` 
	// 0:待揽收 1:已揽收 2:运输中 3:派送中 4:已签收 5:已取消
	SiteID        uint
	DeliveryManID uint
	CreateTime    time.Time
	SignTime      *time.Time
}

// 物流信息表
type Logistics struct {
	gorm.Model
	OrderID     uint
	OrderNo     string
	Status      int    // 状态码，同Order表
	Description string
	Location    string
	Operator    string
	CreateTime  time.Time
}

// 配送信息表
type Delivery struct {
	gorm.Model
	OrderID       uint
	OrderNo       string
	DeliveryManID uint
	SiteID        uint
	Status        int // 0:待配送 1:配送中 2:已完成 3:已退回
	StartTime     *time.Time
	EndTime       *time.Time
	Notes         string
}

// 报销表
type Reimbursement struct {
	gorm.Model
	DeliveryManID uint
	OrderID       uint
	OrderNo       string
	Type          int    // 1:油费 2:过路费 3:维修费 4:其他
	Amount        float64
	Description   string
	Proof         string // 凭证图片路径
	Status        int    // 0:待审核 1:已通过 2:已拒绝
	AdminID       uint
	AdminRemark   string
	ApplyTime     time.Time
	AuditTime     *time.Time
}

// 仓库表
type Warehouse struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Location string
	Capacity int
	Used     int `gorm:"default:0"`
	Manager  string
	Phone    string
	Status   int `gorm:"default:1"` // 1:正常 0:关闭
}

// 仓库物品表
type WarehouseItem struct {
	gorm.Model
	WarehouseID  uint
	OrderID      uint
	OrderNo      string
	LocationCode string // 库位编码
	Quantity     int
	Status       int `gorm:"default:0"` // 0:在库 1:出库
	InTime       time.Time
	OutTime      *time.Time
}

// 结算表
type Settlement struct {
	gorm.Model
	DeliveryManID uint
	StartDate     time.Time
	EndDate       time.Time
	OrderCount    int
	TotalAmount   float64
	DeductAmount  float64
	ActualAmount  float64
	Status        int `gorm:"default:0"` // 0:待结算 1:已结算
	AdminID       uint
	SettleTime    *time.Time
	Notes         string
}
