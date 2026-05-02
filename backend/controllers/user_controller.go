package controllers

import (
	"net/http"
	"time"

	"logisticsManage/models"
	"logisticsManage/utils"

	"github.com/gin-gonic/gin"
)

// 用户注册
func UserRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Phone    string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if models.DB.Where("username = ?", req.Username).First(&existingUser).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Phone:    req.Phone,
		Status:   1,
	}

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "注册失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
	})
}

// 用户登录
func UserLogin(c *gin.Context) {
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

	// 查找用户
	var user models.User
	if err := models.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "账户已被禁用",
		})
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 生成Token
	token, err := utils.GenerateToken(user.ID, user.Username, "user")
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
				"id":        user.ID,
				"username":  user.Username,
				"phone":     user.Phone,
				"email":     user.Email,
				"real_name": user.RealName,
				"address":   user.Address,
			},
		},
	})
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

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
		"data": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"phone":     user.Phone,
			"email":     user.Email,
			"real_name": user.RealName,
			"address":   user.Address,
		},
	})
}

// 更新用户信息
func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		RealName string `json:"real_name"`
		Address  string `json:"address"`
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

	// 更新用户信息
	updates := make(map[string]interface{})
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Address != "" {
		updates["address"] = req.Address
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

// 修改密码
func UpdateUserPassword(c *gin.Context) {
	userID := c.GetUint("user_id")

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

	var user models.User
	if err := models.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	// 验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "旧密码错误",
		})
		return
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	// 更新密码
	if err := models.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
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

// 站点查询
func GetSites(c *gin.Context) {
	var sites []models.Site
	query := models.DB.Where("status = ?", 1)

	// 支持按城市搜索
	city := c.Query("city")
	if city != "" {
		query = query.Where("city LIKE ?", "%"+city+"%")
	}

	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("name LIKE ? OR address LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

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

// 获取站点详情
func GetSiteDetail(c *gin.Context) {
	id := c.Param("id")

	var site models.Site
	if err := models.DB.First(&site, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "站点不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": site,
	})
}

// 生成订单号
func generateOrderNo() string {
	// 格式: 年月日时分秒 + 6位随机数
	now := time.Now()
	return now.Format("20060102150405") + "000001"
}

// 创建订单（寄件）
func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		SenderName    string  `json:"sender_name" binding:"required"`
		SenderPhone   string  `json:"sender_phone" binding:"required"`
		SenderAddress string  `json:"sender_address" binding:"required"`
		ReceiverName  string  `json:"receiver_name" binding:"required"`
		ReceiverPhone string  `json:"receiver_phone" binding:"required"`
		ReceiverAddr  string  `json:"receiver_addr" binding:"required"`
		GoodsType     string  `json:"goods_type"`
		GoodsWeight   float64 `json:"goods_weight"`
		GoodsAmount   int     `json:"goods_amount"`
		InsuranceFee  float64 `json:"insurance_fee"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 计算运费（简单计算：重量*10元/kg + 保险费）
	totalFee := req.GoodsWeight*10 + req.InsuranceFee

	order := models.Order{
		OrderNo:       generateOrderNo(),
		UserID:        userID,
		SenderName:    req.SenderName,
		SenderPhone:   req.SenderPhone,
		SenderAddress: req.SenderAddress,
		ReceiverName:  req.ReceiverName,
		ReceiverPhone: req.ReceiverPhone,
		ReceiverAddr:  req.ReceiverAddr,
		GoodsType:     req.GoodsType,
		GoodsWeight:   req.GoodsWeight,
		GoodsAmount:   req.GoodsAmount,
		InsuranceFee:  req.InsuranceFee,
		TotalFee:      totalFee,
		Status:        0,
		CreateTime:    time.Now(),
	}

	if err := models.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "订单创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "订单创建成功",
		"data": gin.H{
			"order_no": order.OrderNo,
			"total_fee": order.TotalFee,
		},
	})
}

// 获取订单列表
func GetUserOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	var orders []models.Order
	query := models.DB.Where("user_id = ?", userID)

	// 按状态筛选
	status := c.Query("status")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 排序
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
func GetOrderDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderNo := c.Param("order_no")

	var order models.Order
	if err := models.DB.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 获取物流信息
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

// 取消订单
func CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderNo := c.Param("order_no")

	var order models.Order
	if err := models.DB.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 只有待揽收的订单可以取消
	if order.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "当前订单状态不支持取消",
		})
		return
	}

	// 更新订单状态
	order.Status = 5
	if err := models.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "取消失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "订单已取消",
	})
}

// 物流查询（根据订单号）
func QueryLogistics(c *gin.Context) {
	orderNo := c.Query("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请输入订单号",
		})
		return
	}

	// 先查找订单
	var order models.Order
	if err := models.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	// 获取物流信息
	var logistics []models.Logistics
	if err := models.DB.Where("order_id = ?", order.ID).Order("created_at DESC").Find(&logistics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"order":     order,
			"logistics": logistics,
		},
	})
}
