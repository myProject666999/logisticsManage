import request from '@/utils/request'

// 认证相关
export const loginAPI = (data: { username: string; password: string; role: string }) => {
  const urlMap: Record<string, string> = {
    user: '/user/login',
    delivery: '/delivery/login',
    admin: '/admin/login'
  }
  return request.post(urlMap[data.role] || '/user/login', {
    username: data.username,
    password: data.password
  })
}

export const registerAPI = (data: { username: string; password: string; phone?: string }) => {
  return request.post('/user/register', data)
}

// 用户相关
export const getUserInfoAPI = (role: string) => {
  const urlMap: Record<string, string> = {
    user: '/user/info',
    delivery: '/delivery/info',
    admin: '/admin/info'
  }
  return request.get(urlMap[role] || '/user/info')
}

export const updateUserInfoAPI = (role: string, data: any) => {
  const urlMap: Record<string, string> = {
    user: '/user/info',
    delivery: '/delivery/info',
    admin: '/admin/info'
  }
  return request.put(urlMap[role] || '/user/info', data)
}

export const updatePasswordAPI = (role: string, data: { old_password: string; new_password: string }) => {
  const urlMap: Record<string, string> = {
    user: '/user/password',
    delivery: '/delivery/password',
    admin: '/admin/password'
  }
  return request.put(urlMap[role] || '/user/password', data)
}

// 站点相关
export const getSitesAPI = (params?: { city?: string; keyword?: string }) => {
  return request.get('/sites', { params })
}

export const getSiteDetailAPI = (id: number) => {
  return request.get(`/sites/${id}`)
}

// 订单相关（用户）
export const createOrderAPI = (data: any) => {
  return request.post('/user/order', data)
}

export const getUserOrdersAPI = (params?: { status?: string }) => {
  return request.get('/user/orders', { params })
}

export const getUserOrderDetailAPI = (orderNo: string) => {
  return request.get(`/user/order/${orderNo}`)
}

export const cancelOrderAPI = (orderNo: string) => {
  return request.put(`/user/order/${orderNo}/cancel`)
}

// 配送员订单相关
export const getPendingOrdersAPI = () => {
  return request.get('/delivery/orders/pending')
}

export const getDeliveryOrdersAPI = (params?: { status?: string }) => {
  return request.get('/delivery/orders', { params })
}

export const pickupOrderAPI = (orderNo: string) => {
  return request.put(`/delivery/order/${orderNo}/pickup`)
}

export const updateLogisticsAPI = (orderNo: string, data: { status: number; description?: string; location?: string }) => {
  return request.put(`/delivery/logistics/${orderNo}`, data)
}

export const updateDeliveryStatusAPI = (orderNo: string, data: { status: number; notes?: string }) => {
  return request.put(`/delivery/delivery/${orderNo}/status`, data)
}

// 报销相关（配送员）
export const applyReimbursementAPI = (data: any) => {
  return request.post('/delivery/reimbursement', data)
}

export const getReimbursementsAPI = (params?: { status?: string }) => {
  return request.get('/delivery/reimbursements', { params })
}

// 物流查询（公开）
export const queryLogisticsAPI = (orderNo: string) => {
  return request.get('/logistics/query', { params: { order_no: orderNo } })
}

// 管理员 - 用户管理
export const getUsersAPI = (params?: { keyword?: string }) => {
  return request.get('/admin/users', { params })
}

export const getAllUsersAPI = (params?: { keyword?: string }) => {
  return request.get('/admin/users', { params })
}

export const createUserAPI = (data: any) => {
  return request.post('/admin/user', data)
}

export const getUserDetailAPI = (id: number) => {
  return request.get(`/admin/user/${id}`)
}

export const updateUserAPI = (id: number, data: any) => {
  return request.put(`/admin/user/${id}`, data)
}

export const deleteUserAPI = (id: number) => {
  return request.delete(`/admin/user/${id}`)
}

export const updateUserStatusAPI = (id: number, status: number) => {
  return request.put(`/admin/user/${id}/status`, { status })
}

// 管理员 - 配送员管理
export const getDeliveryMenAPI = (params?: { keyword?: string; verify_status?: string }) => {
  return request.get('/admin/delivery-men', { params })
}

export const getAllDeliveryMenAPI = (params?: { keyword?: string; verify_status?: string }) => {
  return request.get('/admin/delivery-men', { params })
}

export const createDeliveryManAPI = (data: any) => {
  return request.post('/admin/delivery-man', data)
}

export const getDeliveryManDetailAPI = (id: number) => {
  return request.get(`/admin/delivery-man/${id}`)
}

export const updateDeliveryManAPI = (id: number, data: any) => {
  return request.put(`/admin/delivery-man/${id}`, data)
}

export const deleteDeliveryManAPI = (id: number) => {
  return request.delete(`/admin/delivery-man/${id}`)
}

export const auditDeliveryManAPI = (id: number, data: { verify_status: number; remark?: string }) => {
  return request.put(`/admin/delivery-man/${id}/audit`, data)
}

export const updateDeliveryManStatusAPI = (id: number, status: number) => {
  return request.put(`/admin/delivery-man/${id}/status`, { status })
}

// 管理员 - 站点管理
export const getAdminSitesAPI = (params?: { keyword?: string }) => {
  return request.get('/admin/sites', { params })
}

export const getAllSitesAPI = (params?: { keyword?: string }) => {
  return request.get('/admin/sites', { params })
}

export const createSiteAPI = (data: any) => {
  return request.post('/admin/site', data)
}

export const updateSiteAPI = (id: number, data: any) => {
  return request.put(`/admin/site/${id}`, data)
}

// 管理员 - 订单管理
export const getAdminOrdersAPI = (params?: { keyword?: string; status?: string }) => {
  return request.get('/admin/orders', { params })
}

export const getAllOrdersAPI = (params?: { keyword?: string; status?: string }) => {
  return request.get('/admin/orders', { params })
}

export const getAdminOrderDetailAPI = (orderNo: string) => {
  return request.get(`/admin/order/${orderNo}`)
}

// 管理员 - 仓库管理
export const getWarehousesAPI = (params?: { keyword?: string }) => {
  return request.get('/admin/warehouses', { params })
}

export const createWarehouseAPI = (data: any) => {
  return request.post('/admin/warehouse', data)
}

export const updateWarehouseAPI = (id: number, data: any) => {
  return request.put(`/admin/warehouse/${id}`, data)
}

export const getWarehouseStatsAPI = () => {
  return request.get('/admin/warehouse/stats')
}

// 管理员 - 物流管理
export const getLogisticsListAPI = (params?: { order_no?: string }) => {
  return request.get('/admin/logistics', { params })
}

export const getLogisticsStatsAPI = () => {
  return request.get('/admin/logistics/stats')
}

// 管理员 - 配送信息管理
export const getDeliveriesAPI = (params?: { status?: string }) => {
  return request.get('/admin/deliveries', { params })
}

// 管理员 - 报销管理
export const getAllReimbursementsAPI = (params?: { status?: string }) => {
  return request.get('/admin/reimbursements', { params })
}

export const auditReimbursementAPI = (id: number, data: { status: number; remark?: string }) => {
  return request.put(`/admin/reimbursement/${id}/audit`, data)
}

export const getReimbursementStatsAPI = () => {
  return request.get('/admin/reimbursement/stats')
}

// 管理员 - 结算管理
export const getSettlementsAPI = (params?: { status?: string }) => {
  return request.get('/admin/settlements', { params })
}

export const getAllSettlementsAPI = (params?: { status?: string }) => {
  return request.get('/admin/settlements', { params })
}

export const createSettlementAPI = (data: any) => {
  return request.post('/admin/settlement', data)
}

export const confirmSettlementAPI = (id: number) => {
  return request.put(`/admin/settlement/${id}/confirm`)
}

export const getSettlementStatsAPI = () => {
  return request.get('/admin/settlement/stats')
}

// 管理员 - 个人信息
export const getAdminProfileAPI = () => {
  return request.get('/admin/info')
}

export const updateAdminProfileAPI = (data: any) => {
  return request.put('/admin/info', data)
}

export const updateAdminPasswordAPI = (data: { old_password: string; new_password: string }) => {
  return request.put('/admin/password', data)
}
