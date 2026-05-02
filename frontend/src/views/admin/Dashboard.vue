<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">用户数</span>
            <el-icon class="card-icon"><User /></el-icon>
          </div>
          <div class="card-value">{{ stats.userCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">配送员数</span>
            <el-icon class="card-icon"><Truck /></el-icon>
          </div>
          <div class="card-value">{{ stats.deliveryManCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">订单总数</span>
            <el-icon class="card-icon"><Document /></el-icon>
          </div>
          <div class="card-value">{{ stats.orderCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">站点数</span>
            <el-icon class="card-icon"><Location /></el-icon>
          </div>
          <div class="card-value">{{ stats.siteCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">待审核报销</span>
            <el-icon class="card-icon"><Money /></el-icon>
          </div>
          <div class="card-value">{{ stats.reimbursePending }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">待结算</span>
            <el-icon class="card-icon"><Coin /></el-icon>
          </div>
          <div class="card-value">{{ stats.settlementPending }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最近订单</span>
          </template>
          <el-table :data="recentOrders" style="width: 100%" max-height="400">
            <el-table-column prop="order_no" label="订单号" width="200" />
            <el-table-column prop="sender_name" label="寄件人" width="100" />
            <el-table-column prop="receiver_name" label="收件人" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-row :gutter="15">
              <el-col :span="12">
                <el-button type="primary" size="large" @click="goTo('/admin/users')">
                  <el-icon><User /></el-icon>
                  用户管理
                </el-button>
              </el-col>
              <el-col :span="12">
                <el-button type="success" size="large" @click="goTo('/admin/delivery-men')">
                  <el-icon><Truck /></el-icon>
                  配送员管理
                </el-button>
              </el-col>
            </el-row>
            <el-row :gutter="15" style="margin-top: 15px">
              <el-col :span="12">
                <el-button type="warning" size="large" @click="goTo('/admin/orders')">
                  <el-icon><Document /></el-icon>
                  订单管理
                </el-button>
              </el-col>
              <el-col :span="12">
                <el-button type="danger" size="large" @click="goTo('/admin/reimbursements')">
                  <el-icon><Money /></el-icon>
                  报销审核
                </el-button>
              </el-col>
            </el-row>
            <el-row :gutter="15" style="margin-top: 15px">
              <el-col :span="12">
                <el-button type="info" size="large" @click="goTo('/admin/sites')">
                  <el-icon><Location /></el-icon>
                  站点管理
                </el-button>
              </el-col>
              <el-col :span="12">
                <el-button type="default" size="large" @click="goTo('/admin/warehouses')">
                  <el-icon><Box /></el-icon>
                  仓库管理
                </el-button>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  getAllUsersAPI, 
  getAllDeliveryMenAPI, 
  getAllOrdersAPI,
  getAllSitesAPI,
  getAllReimbursementsAPI,
  getAllSettlementsAPI
} from '@/api'

const router = useRouter()

const stats = ref({
  userCount: 0,
  deliveryManCount: 0,
  orderCount: 0,
  siteCount: 0,
  reimbursePending: 0,
  settlementPending: 0
})

const recentOrders = ref<any[]>([])

const loadData = async () => {
  try {
    const usersRes: any = await getAllUsersAPI()
    stats.value.userCount = (usersRes.data || []).length
  } catch {}
  
  try {
    const deliveryRes: any = await getAllDeliveryMenAPI()
    stats.value.deliveryManCount = (deliveryRes.data || []).length
  } catch {}
  
  try {
    const ordersRes: any = await getAllOrdersAPI()
    const orders = ordersRes.data || []
    stats.value.orderCount = orders.length
    recentOrders.value = orders.slice(0, 8)
  } catch {}
  
  try {
    const sitesRes: any = await getAllSitesAPI()
    stats.value.siteCount = (sitesRes.data || []).length
  } catch {}
  
  try {
    const reimburseRes: any = await getAllReimbursementsAPI({ status: '0' })
    stats.value.reimbursePending = (reimburseRes.data || []).length
  } catch {}
  
  try {
    const settleRes: any = await getAllSettlementsAPI({ status: '0' })
    stats.value.settlementPending = (settleRes.data || []).length
  } catch {}
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待揽收',
    1: '已揽收',
    2: '运输中',
    3: '派送中',
    4: '已签收',
    5: '已取消'
  }
  return map[status] || '未知'
}

const getStatusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'warning',
    1: 'info',
    2: 'primary',
    3: 'primary',
    4: 'success',
    5: 'danger'
  }
  return map[status] || 'info'
}

const goTo = (path: string) => router.push(path)

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.dashboard-container {
  padding: 0;
}

.card-item {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 14px;
  color: #666;
}

.card-icon {
  font-size: 20px;
  color: #409eff;
}

.card-value {
  font-size: 32px;
  font-weight: bold;
  color: #333;
  margin-top: 15px;
}

.quick-actions .el-button {
  width: 100%;
  height: 50px;
}
</style>
