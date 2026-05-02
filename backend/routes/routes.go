package routes

import (
	"logisticsManage/controllers"
	"logisticsManage/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// API版本分组
	api := r.Group("/api")
	{
		// 公开路由 - 不需要认证
		public := api.Group("")
		{
			// 用户相关
			public.POST("/user/register", controllers.UserRegister)
			public.POST("/user/login", controllers.UserLogin)
			
			// 配送员登录
			public.POST("/delivery/login", controllers.DeliveryLogin)
			
			// 管理员登录
			public.POST("/admin/login", controllers.AdminLogin)
			
			// 公开的物流查询
			public.GET("/logistics/query", controllers.QueryLogistics)
			
			// 公开的站点查询
			public.GET("/sites", controllers.GetSites)
			public.GET("/sites/:id", controllers.GetSiteDetail)
		}

		// ================= 用户路由组 =================
		user := api.Group("/user")
		user.Use(middleware.JWTAuth(), middleware.RoleAuth("user"))
		{
			// 个人信息
			user.GET("/info", controllers.GetUserInfo)
			user.PUT("/info", controllers.UpdateUserInfo)
			user.PUT("/password", controllers.UpdateUserPassword)
			
			// 订单管理
			user.POST("/order", controllers.CreateOrder)
			user.GET("/orders", controllers.GetUserOrders)
			user.GET("/order/:order_no", controllers.GetOrderDetail)
			user.PUT("/order/:order_no/cancel", controllers.CancelOrder)
		}

		// ================= 配送员路由组 =================
		delivery := api.Group("/delivery")
		delivery.Use(middleware.JWTAuth(), middleware.RoleAuth("delivery"))
		{
			// 个人信息
			delivery.GET("/info", controllers.GetDeliveryInfo)
			delivery.PUT("/info", controllers.UpdateDeliveryInfo)
			delivery.PUT("/password", controllers.UpdateDeliveryPassword)
			
			// 订单管理
			delivery.GET("/orders/pending", controllers.GetPendingOrders)
			delivery.GET("/orders", controllers.GetDeliveryOrders)
			
			// 揽收入库
			delivery.PUT("/order/:order_no/pickup", controllers.PickupOrder)
			
			// 物流信息管理
			delivery.PUT("/logistics/:order_no", controllers.UpdateLogistics)
			
			// 配送信息管理
			delivery.PUT("/delivery/:order_no/status", controllers.UpdateDeliveryStatus)
			
			// 报销管理
			delivery.POST("/reimbursement", controllers.ApplyReimbursement)
			delivery.GET("/reimbursements", controllers.GetReimbursements)
		}

		// ================= 管理员路由组 =================
		admin := api.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.RoleAuth("admin"))
		{
			// 个人信息
			admin.GET("/info", controllers.GetAdminInfo)
			admin.PUT("/info", controllers.UpdateAdminInfo)
			admin.PUT("/password", controllers.UpdateAdminPassword)
			
			// 用户管理
			admin.GET("/users", controllers.GetUserList)
			admin.POST("/user", controllers.CreateUser)
			admin.GET("/user/:id", controllers.GetUserDetail)
			admin.PUT("/user/:id", controllers.UpdateUser)
			admin.DELETE("/user/:id", controllers.DeleteUser)
			admin.PUT("/user/:id/status", controllers.UpdateUserStatus)
			
			// 配送员管理
			admin.GET("/delivery-men", controllers.GetDeliveryManList)
			admin.POST("/delivery-man", controllers.CreateDeliveryMan)
			admin.GET("/delivery-man/:id", controllers.GetDeliveryManDetail)
			admin.PUT("/delivery-man/:id", controllers.UpdateDeliveryMan)
			admin.DELETE("/delivery-man/:id", controllers.DeleteDeliveryMan)
			admin.PUT("/delivery-man/:id/audit", controllers.AuditDeliveryMan)
			admin.PUT("/delivery-man/:id/status", controllers.UpdateDeliveryManStatus)
			
			// 站点管理
			admin.GET("/sites", controllers.GetSiteList)
			admin.POST("/site", controllers.CreateSite)
			admin.PUT("/site/:id", controllers.UpdateSite)
			
			// 订单管理
			admin.GET("/orders", controllers.GetOrderList)
			admin.GET("/order/:order_no", controllers.GetAdminOrderDetail)
			
			// 仓库管理
			admin.GET("/warehouses", controllers.GetWarehouseList)
			admin.POST("/warehouse", controllers.CreateWarehouse)
			admin.PUT("/warehouse/:id", controllers.UpdateWarehouse)
			
			// 仓库统计
			admin.GET("/warehouse/stats", controllers.GetWarehouseStats)
			
			// 物流管理
			admin.GET("/logistics", controllers.GetLogisticsList)
			
			// 物流统计
			admin.GET("/logistics/stats", controllers.GetLogisticsStats)
			
			// 配送信息管理
			admin.GET("/deliveries", controllers.GetDeliveryList)
			
			// 报销管理
			admin.GET("/reimbursements", controllers.GetAllReimbursements)
			admin.PUT("/reimbursement/:id/audit", controllers.AuditReimbursement)
			
			// 报销统计
			admin.GET("/reimbursement/stats", controllers.GetReimbursementStats)
			
			// 结算管理
			admin.GET("/settlements", controllers.GetSettlementList)
			admin.POST("/settlement", controllers.CreateSettlement)
			admin.PUT("/settlement/:id/confirm", controllers.ConfirmSettlement)
			
			// 结算统计
			admin.GET("/settlement/stats", controllers.GetSettlementStats)
		}
	}
}
