import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { ElMessage } from 'element-plus'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  // 用户端路由
  {
    path: '/user',
    name: 'UserLayout',
    component: () => import('@/layouts/UserLayout.vue'),
    meta: { requiresAuth: true, role: 'user' },
    children: [
      {
        path: '',
        redirect: '/user/dashboard'
      },
      {
        path: 'dashboard',
        name: 'UserDashboard',
        component: () => import('@/views/user/Dashboard.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'sites',
        name: 'UserSites',
        component: () => import('@/views/user/Sites.vue'),
        meta: { title: '站点查询' }
      },
      {
        path: 'order/create',
        name: 'CreateOrder',
        component: () => import('@/views/user/CreateOrder.vue'),
        meta: { title: '寄件' }
      },
      {
        path: 'orders',
        name: 'UserOrders',
        component: () => import('@/views/user/Orders.vue'),
        meta: { title: '订单管理' }
      },
      {
        path: 'logistics',
        name: 'UserLogistics',
        component: () => import('@/views/user/Logistics.vue'),
        meta: { title: '物流查询' }
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/user/Profile.vue'),
        meta: { title: '个人信息' }
      },
      {
        path: 'password',
        name: 'UserPassword',
        component: () => import('@/views/user/Password.vue'),
        meta: { title: '密码修改' }
      }
    ]
  },
  // 配送员端路由
  {
    path: '/delivery',
    name: 'DeliveryLayout',
    component: () => import('@/layouts/DeliveryLayout.vue'),
    meta: { requiresAuth: true, role: 'delivery' },
    children: [
      {
        path: '',
        redirect: '/delivery/dashboard'
      },
      {
        path: 'dashboard',
        name: 'DeliveryDashboard',
        component: () => import('@/views/delivery/Dashboard.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'orders/pending',
        name: 'PendingOrders',
        component: () => import('@/views/delivery/PendingOrders.vue'),
        meta: { title: '待揽收订单' }
      },
      {
        path: 'orders',
        name: 'DeliveryOrders',
        component: () => import('@/views/delivery/Orders.vue'),
        meta: { title: '订单管理' }
      },
      {
        path: 'logistics',
        name: 'DeliveryLogistics',
        component: () => import('@/views/delivery/Logistics.vue'),
        meta: { title: '物流信息管理' }
      },
      {
        path: 'delivery-info',
        name: 'DeliveryInfo',
        component: () => import('@/views/delivery/DeliveryInfo.vue'),
        meta: { title: '配送信息管理' }
      },
      {
        path: 'reimbursements',
        name: 'DeliveryReimbursements',
        component: () => import('@/views/delivery/Reimbursements.vue'),
        meta: { title: '报销管理' }
      },
      {
        path: 'profile',
        name: 'DeliveryProfile',
        component: () => import('@/views/delivery/Profile.vue'),
        meta: { title: '个人信息' }
      },
      {
        path: 'password',
        name: 'DeliveryPassword',
        component: () => import('@/views/delivery/Password.vue'),
        meta: { title: '密码修改' }
      }
    ]
  },
  // 管理员端路由
  {
    path: '/admin',
    name: 'AdminLayout',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true, role: 'admin' },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('@/views/admin/Dashboard.vue'),
        meta: { title: '首页' }
      },
      // 用户管理
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('@/views/admin/Users.vue'),
        meta: { title: '用户管理' }
      },
      // 配送员管理
      {
        path: 'delivery-men',
        name: 'AdminDeliveryMen',
        component: () => import('@/views/admin/DeliveryMen.vue'),
        meta: { title: '配送员管理' }
      },
      // 站点管理
      {
        path: 'sites',
        name: 'AdminSites',
        component: () => import('@/views/admin/Sites.vue'),
        meta: { title: '站点管理' }
      },
      // 订单管理
      {
        path: 'orders',
        name: 'AdminOrders',
        component: () => import('@/views/admin/Orders.vue'),
        meta: { title: '订单管理' }
      },
      // 仓库管理
      {
        path: 'warehouses',
        name: 'AdminWarehouses',
        component: () => import('@/views/admin/Warehouses.vue'),
        meta: { title: '仓库管理' }
      },
      {
        path: 'warehouse-stats',
        name: 'WarehouseStats',
        component: () => import('@/views/admin/WarehouseStats.vue'),
        meta: { title: '仓库统计' }
      },
      // 物流管理
      {
        path: 'logistics',
        name: 'AdminLogistics',
        component: () => import('@/views/admin/Logistics.vue'),
        meta: { title: '物流管理' }
      },
      {
        path: 'logistics-stats',
        name: 'LogisticsStats',
        component: () => import('@/views/admin/LogisticsStats.vue'),
        meta: { title: '物流统计' }
      },
      // 配送信息管理
      {
        path: 'deliveries',
        name: 'AdminDeliveries',
        component: () => import('@/views/admin/Deliveries.vue'),
        meta: { title: '配送信息管理' }
      },
      // 报销管理
      {
        path: 'reimbursements',
        name: 'AdminReimbursements',
        component: () => import('@/views/admin/Reimbursements.vue'),
        meta: { title: '报销管理' }
      },
      {
        path: 'reimbursement-stats',
        name: 'ReimbursementStats',
        component: () => import('@/views/admin/ReimbursementStats.vue'),
        meta: { title: '报销统计' }
      },
      // 结算管理
      {
        path: 'settlements',
        name: 'AdminSettlements',
        component: () => import('@/views/admin/Settlements.vue'),
        meta: { title: '结算管理' }
      },
      {
        path: 'settlement-stats',
        name: 'SettlementStats',
        component: () => import('@/views/admin/SettlementStats.vue'),
        meta: { title: '结算统计' }
      },
      // 个人信息和密码
      {
        path: 'profile',
        name: 'AdminProfile',
        component: () => import('@/views/admin/Profile.vue'),
        meta: { title: '个人信息' }
      },
      {
        path: 'password',
        name: 'AdminPassword',
        component: () => import('@/views/admin/Password.vue'),
        meta: { title: '密码修改' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  const role = localStorage.getItem('role')

  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 物流管理系统`
  }

  // 不需要登录的页面
  if (['/login', '/register'].includes(to.path)) {
    if (token && role) {
      // 已登录用户跳转到对应首页
      const redirectMap: Record<string, string> = {
        user: '/user/dashboard',
        delivery: '/delivery/dashboard',
        admin: '/admin/dashboard'
      }
      next(redirectMap[role] || '/login')
      return
    }
    next()
    return
  }

  // 需要登录的页面
  if (to.meta.requiresAuth) {
    if (!token) {
      next('/login')
      return
    }

    // 检查角色权限
    if (to.meta.role && to.meta.role !== role) {
      ElMessage.error('没有权限访问此页面')
      const redirectMap: Record<string, string> = {
        user: '/user/dashboard',
        delivery: '/delivery/dashboard',
        admin: '/admin/dashboard'
      }
      next(redirectMap[role || 'user'] || '/login')
      return
    }

    next()
    return
  }

  next()
})

export default router
