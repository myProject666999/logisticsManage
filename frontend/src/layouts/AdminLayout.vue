<template>
  <el-container class="layout-container">
    <el-aside width="220px" class="layout-aside">
      <div class="logo">
        <el-icon><Box /></el-icon>
        <span>物流管理系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        router
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><HomeFilled /></el-icon>
          <span>首页</span>
        </el-menu-item>
        <el-sub-menu index="user-manage">
          <template #title>
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </template>
          <el-menu-item index="/admin/users">用户列表</el-menu-item>
          <el-menu-item index="/admin/delivery-men">配送员管理</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/admin/sites">
          <el-icon><OfficeBuilding /></el-icon>
          <span>站点管理</span>
        </el-menu-item>
        <el-menu-item index="/admin/orders">
          <el-icon><Document /></el-icon>
          <span>订单管理</span>
        </el-menu-item>
        <el-sub-menu index="warehouse">
          <template #title>
            <el-icon><Warehouse /></el-icon>
            <span>仓库管理</span>
          </template>
          <el-menu-item index="/admin/warehouses">仓库列表</el-menu-item>
          <el-menu-item index="/admin/warehouse-stats">仓库统计</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="logistics">
          <template #title>
            <el-icon><MapLocation /></el-icon>
            <span>物流管理</span>
          </template>
          <el-menu-item index="/admin/logistics">物流列表</el-menu-item>
          <el-menu-item index="/admin/logistics-stats">物流统计</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/admin/deliveries">
          <el-icon><Truck /></el-icon>
          <span>配送信息</span>
        </el-menu-item>
        <el-sub-menu index="reimbursement">
          <template #title>
            <el-icon><Money /></el-icon>
            <span>报销管理</span>
          </template>
          <el-menu-item index="/admin/reimbursements">报销列表</el-menu-item>
          <el-menu-item index="/admin/reimbursement-stats">报销统计</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="settlement">
          <template #title>
            <el-icon><Wallet /></el-icon>
            <span>结算管理</span>
          </template>
          <el-menu-item index="/admin/settlements">结算列表</el-menu-item>
          <el-menu-item index="/admin/settlement-stats">结算统计</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="user">
          <template #title>
            <el-icon><User /></el-icon>
            <span>个人中心</span>
          </template>
          <el-menu-item index="/admin/profile">个人信息</el-menu-item>
          <el-menu-item index="/admin/password">密码修改</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header class="layout-header">
        <div class="header-title">
          <span>管理员端</span>
        </div>
        <div class="header-user">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-icon><UserFilled /></el-icon>
              {{ userInfo?.username || '管理员' }}
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="password">修改密码</el-dropdown-item>
                <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="layout-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)
const userInfo = computed(() => userStore.userInfo)

const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/admin/profile')
      break
    case 'password':
      router.push('/admin/password')
      break
    case 'logout':
      userStore.logout()
      router.push('/login')
      break
  }
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}

.layout-aside {
  background-color: #304156;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  border-bottom: 1px solid #1f2d3d;
}

.logo .el-icon {
  margin-right: 10px;
  font-size: 24px;
}

.layout-header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-title {
  font-size: 18px;
  font-weight: 500;
  color: #333;
}

.header-user {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  color: #606266;
}

.user-info .el-icon {
  margin-right: 5px;
}

.layout-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>
