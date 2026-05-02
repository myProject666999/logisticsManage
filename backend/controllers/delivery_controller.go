package controllers

import (
	"net/http"
	"time"

	"logisticsManage/models"
	"logisticsManage/utils"

	"github.com/gin-gonic/gin"
)

// 配送员登录
func DeliveryLogin(c *gin.Context) {
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

	// 查找配送员
	var deliveryMan models.DeliveryMan
	if err := models.DB.Where("username = ?", req.Username).First(&deliveryMan).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查状态
	if deliveryMan.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "账户已被禁用",
		})
		return
	}

	if deliveryMan.VerifyStatus != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "账户未通过审核",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, deliveryMan.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(deliveryMan.ID, deliveryMan.Username, "delivery")
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
				"id":         deliveryMan.ID,
				"username":   deliveryMan.Username,
				"phone":      deliveryMan.Phone,
				"real_name":  deliveryMan.RealName,
				"site_id":    deliveryMan.SiteID,
			},
		},
	})
}

// 获取配送员信息
func GetDeliveryInfo(c *gin.Context) {
	deliveryID := c.GetUint("user_id")

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
		"data": gin.H{
			"id":           deliveryMan.ID,
			"username":     deliveryMan.Username,
			"phone":        deliveryMan.Phone,
			"real_name":    deliveryMan.RealName,
			"id_card":      deliveryMan.IDCard,
			"site_id":      deliveryMan.SiteID,
			"verify_status": deliveryMan.VerifyStatus,
		},
	})
}

// 更新配送员信息
func UpdateDeliveryInfo(c *gin.Context) {
	deliveryID := c.GetUint("user_id")

	var req struct {
		Phone    string `json:"phone"`
		RealName string `json:"real_name"`
		IDCard   string `json:"id_card"`
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
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.IDCard != "" {
		updates["id_card"] = req.IDCard
		updates["verify_status"] = 0 // 重新审核
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

// 修改配送员密码
func UpdateDeliveryPassword(c *gin.Context) {
	deliveryID := c.GetUint("user_id")

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

	var deliveryMan models.DeliveryMan
	if err := models.DB.First(&deliveryMan, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "配送员不存在",
		})
		return
	}

	if !utils.CheckPassword(req.OldPassword, deliveryMan.Password) {
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

	if err := models.DB.Model(&deliveryMan).Update("password", hashedPassword).Error; err != nil {
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

// 获取待揽收订单列表
func GetPendingOrders(c *gin.Context) {
	deliveryID := c.GetUint("user_id")

	var deliveryMan models.DeliveryMan
	if err := models.DB.First(&deliveryMan, deliveryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "配送员不存在",
		})
		return
	}

	// 获取本站点的待揽收订单
	var orders []models.Order
	query := models.DB.Where("status = ?", 0) // 待揽收

	// 如果配送员分配了站点，只看该站点的订单
	if deliveryMan.SiteID > 0 {
		query = query.Where("site_id = ?", deliveryMan.SiteID)
	}

	// 或者可以根据地址分配，这里简化处理
	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
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

// 获取配送员的订单列表
func GetDeliveryOrders(c *gin.Context) {
	deliveryID := c.GetUint("user_id")

	var orders []models.Order
	query := models.DB.Where("delivery_man_id = ?", deliveryID)

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

// 揽收入库
func PickupOrder(c *gin.Context) {
	deliveryID := c.GetUint("user_id")
	orderNo := c.Param("order_no")

	var order models.Order
	if err := models.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	if order.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "订单状态不支持揽收",
		})
		return
	}

	// 更新订单状态
	order.Status = 1
	order.DeliveryManID = deliveryID
	if err := models.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "揽收失败",
		})
		return
	}

	// 添加物流记录
	logistics := models.Logistics{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		Status:      1,
		Description: "快递员已揽收",
		Operator:    "配送员",
		CreateTime:  time.Now(),
	}
	models.DB.Create(&logistics)

	// 创建配送记录
	delivery := models.Delivery{
		OrderID:       order.ID,
		OrderNo:       order.OrderNo,
		DeliveryManID: deliveryID,
		SiteID:        order.SiteID,
		Status:        0,
	}
	models.DB.Create(&delivery)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "揽收成功",
	})
}

// 更新物流信息
func UpdateLogistics(c *gin.Context) {
	deliveryID := c.GetUint("user_id")
	orderNo := c.Param("order_no")

	var req struct {
		Status      int    `json:"status" binding:"required"`
		Description string `json:"description"`
		Location    string `json:"location"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var order models.Order
	if err := models.DB.Where("order_no = ? AND delivery_man_id = ?", orderNo, deliveryID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 更新订单状态
	order.Status = req.Status
	if req.Status == 4 { // 已签收
		now := time.Now()
		order.SignTime = &now
	}
	if err := models.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	// 添加物流记录
	logistics := models.Logistics{
		OrderID:     order.ID,
		OrderNo:     order.OrderNo,
		Status:      req.Status,
		Description: req.Description,
		Location:    req.Location,
		Operator:    "配送员",
		CreateTime:  time.Now(),
	}
	models.DB.Create(&logistics)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "物流信息更新成功",
	})
}

// 更新配送信息
func UpdateDeliveryStatus(c *gin.Context) {
	deliveryID := c.GetUint("user_id")
	orderNo := c.Param("order_no")

	var req struct {
		Status int    `json:"status" binding:"required"`
		Notes  string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var delivery models.Delivery
	if err := models.DB.Where("order_no = ? AND delivery_man_id = ?", orderNo, deliveryID).First(&delivery).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "配送记录不存在",
		})
		return
	}

	now := time.Now()
	if req.Status == 1 { // 开始配送
		delivery.StartTime = &now
	} else if req.Status == 2 || req.Status == 3 { // 已完成或已退回
		delivery.EndTime = &now
	}

	delivery.Status = req.Status
	delivery.Notes = req.Notes

	if err := models.DB.Save(&delivery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "配送状态更新成功",
	})
}

// 申请报销
func ApplyReimbursement(c *gin.Context) {
	deliveryID := c.GetUint("user_id")

	var req struct {
		OrderNo     string  `json:"order_no"`
		Type        int     `json:"type" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description"`
		Proof       string  `json:"proof"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var orderID uint
	if req.OrderNo != "" {
		var order models.Order
		if err := models.DB.Where("order_no = ?", req.OrderNo).First(&order).Error; err == nil {
			orderID = order.ID
		}
	}

	reimbursement := models.Reimbursement{
		DeliveryManID: deliveryID,
		OrderID:       orderID,
		OrderNo:       req.OrderNo,
		Type:          req.Type,
		Amount:        req.Amount,
		Description:   req.Description,
		Proof:         req.Proof,
		Status:        0,
		ApplyTime:     time.Now(),
	}

	if err := models.DB.Create(&reimbursement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "申请失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "报销申请提交成功",
	})
}

// 获取报销列表
func GetReimbursements(c *gin.Context) {
	deliveryID := c.GetUint("user_id")

	var reimbursements []models.Reimbursement
	query := models.DB.Where("delivery_man_id = ?", deliveryID)

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
