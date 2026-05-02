package controllers

import (
	"net/http"
	"time"

	"logisticsManage/models"
	"logisticsManage/utils"

	"github.com/gin-gonic/gin"
)

// 管理员登录
func AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var admin models.Admin
	if err := models.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	if admin.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "账户已被禁用",
		})
		return
	}

	if !utils.CheckPassword(req.Password, admin.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Token生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       admin.ID,
				"username": admin.Username,
				"real_name": admin.RealName,
				"phone":    admin.Phone,
				"role":     admin.Role,
			},
		},
	})
}

// 获取管理员信息
func GetAdminInfo(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var admin models.Admin
	if err := models.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"real_name": admin.RealName,
			"phone":    admin.Phone,
			"role":     admin.Role,
		},
	})
}

// 更新管理员信息
func UpdateAdminInfo(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var req struct {
		Phone    string `json:"phone"`
		RealName string `json:"real_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var admin models.Admin
	if err := models.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}

	if err := models.DB.Model(&admin).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// 修改管理员密码
func UpdateAdminPassword(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var admin models.Admin
	if err := models.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	if !utils.CheckPassword(req.OldPassword, admin.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "旧密码错误",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	if err := models.DB.Model(&admin).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码修改成功",
	})
}

// ================= 用户管理 =================

// 获取用户列表
func GetUserList(c *gin.Context) {
	var users []models.User
	query := models.DB

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": users,
	})
}

// 更新用户状态
func UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	user.Status = req.Status
	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "状态更新成功",
	})
}

// ================= 配送员管理 =================

// 获取配送员列表
func GetDeliveryManList(c *gin.Context) {
	var deliveryMen []models.DeliveryMan
	query := models.DB

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("username LIKE ? OR phone LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	verifyStatus := c.Query("verify_status")
	if verifyStatus != "" {
		query = query.Where("verify_status = ?", verifyStatus)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&deliveryMen).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": deliveryMen,
	})
}

// 审核配送员
func AuditDeliveryMan(c *gin.Context) {
	deliveryID := c.Param("id")

	var req struct {
		VerifyStatus int    `json:"verify_status" binding:"required"`
		Remark       string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var deliveryMan models.DeliveryMan
	if err := models.DB.First(&deliveryMan, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "配送员不存在",
		})
		return
	}

	deliveryMan.VerifyStatus = req.VerifyStatus
	if err := models.DB.Save(&deliveryMan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "审核失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "审核完成",
	})
}

// 更新配送员状态
func UpdateDeliveryManStatus(c *gin.Context) {
	deliveryID := c.Param("id")

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var deliveryMan models.DeliveryMan
	if err := models.DB.First(&deliveryMan, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "配送员不存在",
		})
		return
	}

	deliveryMan.Status = req.Status
	if err := models.DB.Save(&deliveryMan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "状态更新成功",
	})
}

// ================= 站点管理 =================

// 获取站点列表
func GetSiteList(c *gin.Context) {
	var sites []models.Site
	query := models.DB

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("name LIKE ? OR city LIKE ? OR address LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&sites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": sites,
	})
}

// 创建站点
func CreateSite(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Province string `json:"province"`
		City     string `json:"city"`
		District string `json:"district"`
		Address  string `json:"address" binding:"required"`
		Phone    string `json:"phone"`
		Manager  string `json:"manager"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	site := models.Site{
		Name:     req.Name,
		Province: req.Province,
		City:     req.City,
		District: req.District,
		Address:  req.Address,
		Phone:    req.Phone,
		Manager:  req.Manager,
		Status:   1,
	}

	if err := models.DB.Create(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "站点创建成功",
	})
}

// 更新站点
func UpdateSite(c *gin.Context) {
	siteID := c.Param("id")

	var req struct {
		Name     string `json:"name"`
		Province string `json:"province"`
		City     string `json:"city"`
		District string `json:"district"`
		Address  string `json:"address"`
		Phone    string `json:"phone"`
		Manager  string `json:"manager"`
		Status   *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var site models.Site
	if err := models.DB.First(&site, siteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "站点不存在",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Province != "" {
		updates["province"] = req.Province
	}
	if req.City != "" {
		updates["city"] = req.City
	}
	if req.District != "" {
		updates["district"] = req.District
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Manager != "" {
		updates["manager"] = req.Manager
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := models.DB.Model(&site).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "站点更新成功",
	})
}

// ================= 订单管理 =================

// 获取订单列表
func GetOrderList(c *gin.Context) {
	var orders []models.Order
	query := models.DB

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("order_no LIKE ? OR sender_name LIKE ? OR receiver_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": orders,
	})
}

// 获取订单详情
func GetAdminOrderDetail(c *gin.Context) {
	orderNo := c.Param("order_no")

	var order models.Order
	if err := models.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	var logistics []models.Logistics
	models.DB.Where("order_id = ?", order.ID).Order("created_at DESC").Find(&logistics)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"order":     order,
			"logistics": logistics,
		},
	})
}

// ================= 仓库管理 =================

// 获取仓库列表
func GetWarehouseList(c *gin.Context) {
	var warehouses []models.Warehouse
	query := models.DB

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("name LIKE ? OR location LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&warehouses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": warehouses,
	})
}

// 创建仓库
func CreateWarehouse(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Location string `json:"location"`
		Capacity int    `json:"capacity"`
		Manager  string `json:"manager"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	warehouse := models.Warehouse{
		Name:     req.Name,
		Location: req.Location,
		Capacity: req.Capacity,
		Manager:  req.Manager,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := models.DB.Create(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "仓库创建成功",
	})
}

// 更新仓库
func UpdateWarehouse(c *gin.Context) {
	warehouseID := c.Param("id")

	var req struct {
		Name     string `json:"name"`
		Location string `json:"location"`
		Capacity int    `json:"capacity"`
		Manager  string `json:"manager"`
		Phone    string `json:"phone"`
		Status   *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var warehouse models.Warehouse
	if err := models.DB.First(&warehouse, warehouseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "仓库不存在",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Capacity > 0 {
		updates["capacity"] = req.Capacity
	}
	if req.Manager != "" {
		updates["manager"] = req.Manager
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := models.DB.Model(&warehouse).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "仓库更新成功",
	})
}

// 仓库统计
func GetWarehouseStats(c *gin.Context) {
	var totalCount int
	var totalCapacity int
	var totalUsed int

	models.DB.Model(&models.Warehouse{}).Count(&totalCount)
	models.DB.Model(&models.Warehouse{}).Select("COALESCE(SUM(capacity), 0)").Scan(&totalCapacity)
	models.DB.Model(&models.Warehouse{}).Select("COALESCE(SUM(used), 0)").Scan(&totalUsed)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total_count":   totalCount,
			"total_capacity": totalCapacity,
			"total_used":    totalUsed,
			"utilization":   float64(totalUsed) / float64(max(totalCapacity, 1)) * 100,
		},
	})
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ================= 物流管理 =================

// 获取物流列表
func GetLogisticsList(c *gin.Context) {
	var logistics []models.Logistics
	query := models.DB

	orderNo := c.Query("order_no")
	if orderNo != "" {
		query = query.Where("order_no = ?", orderNo)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&logistics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": logistics,
	})
}

// 物流统计
func GetLogisticsStats(c *gin.Context) {
	var total int
	var pending int    // 待揽收
	var inTransit int // 运输中/派送中
	var completed int // 已签收
	var cancelled int // 已取消

	models.DB.Model(&models.Order{}).Count(&total)
	models.DB.Model(&models.Order{}).Where("status = ?", 0).Count(&pending)
	models.DB.Model(&models.Order{}).Where("status IN (?)", []int{1, 2, 3}).Count(&inTransit)
	models.DB.Model(&models.Order{}).Where("status = ?", 4).Count(&completed)
	models.DB.Model(&models.Order{}).Where("status = ?", 5).Count(&cancelled)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total":      total,
			"pending":    pending,
			"in_transit": inTransit,
			"completed":  completed,
			"cancelled":  cancelled,
		},
	})
}

// ================= 配送信息管理 =================

// 获取配送列表
func GetDeliveryList(c *gin.Context) {
	var deliveries []models.Delivery
	query := models.DB

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&deliveries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": deliveries,
	})
}

// ================= 报销管理 =================

// 获取报销列表
func GetAllReimbursements(c *gin.Context) {
	var reimbursements []models.Reimbursement
	query := models.DB

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("apply_time DESC")

	if err := query.Find(&reimbursements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": reimbursements,
	})
}

// 审核报销
func AuditReimbursement(c *gin.Context) {
	adminID := c.GetUint("user_id")
	reimbursementID := c.Param("id")

	var req struct {
		Status int    `json:"status" binding:"required"`
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var reimbursement models.Reimbursement
	if err := models.DB.First(&reimbursement, reimbursementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "报销申请不存在",
		})
		return
	}

	now := time.Now()
	reimbursement.Status = req.Status
	reimbursement.AdminID = adminID
	reimbursement.AdminRemark = req.Remark
	reimbursement.AuditTime = &now

	if err := models.DB.Save(&reimbursement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "审核失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "审核完成",
	})
}

// 报销统计
func GetReimbursementStats(c *gin.Context) {
	var totalCount int
	var pendingCount int
	var approvedCount int
	var rejectedCount int
	var totalAmount float64
	var approvedAmount float64

	models.DB.Model(&models.Reimbursement{}).Count(&totalCount)
	models.DB.Model(&models.Reimbursement{}).Where("status = ?", 0).Count(&pendingCount)
	models.DB.Model(&models.Reimbursement{}).Where("status = ?", 1).Count(&approvedCount)
	models.DB.Model(&models.Reimbursement{}).Where("status = ?", 2).Count(&rejectedCount)
	models.DB.Model(&models.Reimbursement{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)
	models.DB.Model(&models.Reimbursement{}).Where("status = ?", 1).Select("COALESCE(SUM(amount), 0)").Scan(&approvedAmount)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total_count":      totalCount,
			"pending_count":    pendingCount,
			"approved_count":   approvedCount,
			"rejected_count":   rejectedCount,
			"total_amount":     totalAmount,
			"approved_amount":  approvedAmount,
		},
	})
}

// ================= 结算管理 =================

// 获取结算列表
func GetSettlementList(c *gin.Context) {
	var settlements []models.Settlement
	query := models.DB

	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&settlements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": settlements,
	})
}

// 创建结算
func CreateSettlement(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var req struct {
		DeliveryManID uint      `json:"delivery_man_id" binding:"required"`
		StartDate     time.Time `json:"start_date" binding:"required"`
		EndDate       time.Time `json:"end_date" binding:"required"`
		Notes         string    `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 计算该配送员在时间段内的订单
	var orderCount int
	var totalAmount float64

	models.DB.Model(&models.Order{}).Where("delivery_man_id = ? AND created_at BETWEEN ? AND ? AND status = ?", 
		req.DeliveryManID, req.StartDate, req.EndDate, 4).Count(&orderCount)
	
	// 简单计算：每单5元
	totalAmount = float64(orderCount) * 5

	settlement := models.Settlement{
		DeliveryManID: req.DeliveryManID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		OrderCount:    orderCount,
		TotalAmount:   totalAmount,
		DeductAmount:  0,
		ActualAmount:  totalAmount,
		Status:        0,
		AdminID:       adminID,
		Notes:         req.Notes,
	}

	if err := models.DB.Create(&settlement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "结算创建成功",
	})
}

// 确认结算
func ConfirmSettlement(c *gin.Context) {
	adminID := c.GetUint("user_id")
	settlementID := c.Param("id")

	var settlement models.Settlement
	if err := models.DB.First(&settlement, settlementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "结算不存在",
		})
		return
	}

	now := time.Now()
	settlement.Status = 1
	settlement.SettleTime = &now
	settlement.AdminID = adminID

	if err := models.DB.Save(&settlement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "结算失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "结算完成",
	})
}

// 结算统计
func GetSettlementStats(c *gin.Context) {
	var totalCount int
	var pendingCount int
	var settledCount int
	var totalAmount float64
	var settledAmount float64

	models.DB.Model(&models.Settlement{}).Count(&totalCount)
	models.DB.Model(&models.Settlement{}).Where("status = ?", 0).Count(&pendingCount)
	models.DB.Model(&models.Settlement{}).Where("status = ?", 1).Count(&settledCount)
	models.DB.Model(&models.Settlement{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalAmount)
	models.DB.Model(&models.Settlement{}).Where("status = ?", 1).Select("COALESCE(SUM(actual_amount), 0)").Scan(&settledAmount)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"total_count":     totalCount,
			"pending_count":   pendingCount,
			"settled_count":   settledCount,
			"total_amount":    totalAmount,
			"settled_amount":  settledAmount,
		},
	})
}

// ================= 补充的CRUD方法 =================

// 创建用户（管理员）
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		RealName string `json:"real_name"`
		IDCard   string `json:"id_card"`
		Address  string `json:"address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existingUser models.User
	if models.DB.Where("username = ?", req.Username).First(&existingUser).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Email:    req.Email,
		RealName: req.RealName,
		Address:  req.Address,
		Status:   1,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": user,
	})
}

// 获取用户详情（管理员）
func GetUserDetail(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
	})
}

// 更新用户（管理员）
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		RealName string `json:"real_name"`
		IDCard   string `json:"id_card"`
		Address  string `json:"address"`
		Status   *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.IDCard != "" {
		updates["id_card"] = req.IDCard
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "密码加密失败",
			})
			return
		}
		updates["password"] = hashedPassword
	}

	if err := models.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// 删除用户（管理员）
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	if err := models.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// 创建配送员（管理员）
func CreateDeliveryMan(c *gin.Context) {
	var req struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		Phone      string `json:"phone"`
		RealName   string `json:"real_name"`
		IDCard     string `json:"id_card"`
		VehicleNo  string `json:"vehicle_no"`
		VehicleType *int   `json:"vehicle_type"`
		SiteID     *uint  `json:"site_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existingDeliveryMan models.DeliveryMan
	if models.DB.Where("username = ?", req.Username).First(&existingDeliveryMan).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	deliveryMan := models.DeliveryMan{
		Username:     req.Username,
		Password:     hashedPassword,
		Phone:        req.Phone,
		RealName:     req.RealName,
		IDCard:       req.IDCard,
		VehicleNo:    req.VehicleNo,
		Status:       1,
		VerifyStatus: 1,
	}

	if req.VehicleType != nil {
		deliveryMan.VehicleType = *req.VehicleType
	}
	if req.SiteID != nil {
		deliveryMan.SiteID = *req.SiteID
	}

	if err := models.DB.Create(&deliveryMan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": deliveryMan,
	})
}

// 获取配送员详情（管理员）
func GetDeliveryManDetail(c *gin.Context) {
	deliveryID := c.Param("id")

	var deliveryMan models.DeliveryMan
	if err := models.DB.First(&deliveryMan, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "配送员不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": deliveryMan,
	})
}

// 更新配送员（管理员）
func UpdateDeliveryMan(c *gin.Context) {
	deliveryID := c.Param("id")

	var req struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		Phone      string `json:"phone"`
		RealName   string `json:"real_name"`
		IDCard     string `json:"id_card"`
		VehicleNo  string `json:"vehicle_no"`
		VehicleType *int   `json:"vehicle_type"`
		SiteID     *uint  `json:"site_id"`
		Status     *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var deliveryMan models.DeliveryMan
	if err := models.DB.First(&deliveryMan, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "配送员不存在",
		})
		return
	}

	updates := make(map[string]interface{})
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.IDCard != "" {
		updates["id_card"] = req.IDCard
	}
	if req.VehicleNo != "" {
		updates["vehicle_no"] = req.VehicleNo
	}
	if req.VehicleType != nil {
		updates["vehicle_type"] = *req.VehicleType
	}
	if req.SiteID != nil {
		updates["site_id"] = *req.SiteID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "密码加密失败",
			})
			return
		}
		updates["password"] = hashedPassword
	}

	if err := models.DB.Model(&deliveryMan).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// 删除配送员（管理员）
func DeleteDeliveryMan(c *gin.Context) {
	deliveryID := c.Param("id")

	var deliveryMan models.DeliveryMan
	if err := models.DB.First(&deliveryMan, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "配送员不存在",
		})
		return
	}

	if err := models.DB.Delete(&deliveryMan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
