<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">待处理订单</span>
            <el-icon class="card-icon"><Clock /></el-icon>
          </div>
          <div class="card-value">{{ stats.pending }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">运输中</span>
            <el-icon class="card-icon"><Truck /></el-icon>
          </div>
          <div class="card-value">{{ stats.inTransit }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">已完成</span>
            <el-icon class="card-icon"><CircleCheck /></el-icon>
          </div>
          <div class="card-value">{{ stats.completed }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">总订单数</span>
            <el-icon class="card-icon"><Document /></el-icon>
          </div>
          <div class="card-value">{{ stats.total }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>最近订单</span>
          </template>
          <el-table :data="recentOrders" style="width: 100%">
            <el-table-column prop="order_no" label="订单号" width="200" />
            <el-table-column prop="receiver_name" label="收件人" width="120" />
            <el-table-column prop="receiver_addr" label="收件地址" />
            <el-table-column prop="total_fee" label="运费(元)" width="100">
              <template #default="scope">
                ¥{{ scope.row.total_fee }}
              </template>
            </el-table-column>
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
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="handleCreateOrder">
              <el-icon><Plus /></el-icon>
              寄件
            </el-button>
            <el-button type="success" size="large" @click="handleViewSites">
              <el-icon><OfficeBuilding /></el-icon>
              查询站点
            </el-button>
            <el-button type="warning" size="large" @click="handleViewOrders">
              <el-icon><Document /></el-icon>
              我的订单
            </el-button>
            <el-button type="info" size="large" @click="handleQueryLogistics">
              <el-icon><Search /></el-icon>
              物流查询
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getUserOrdersAPI } from '@/api'

const router = useRouter()

const stats = ref({
  pending: 0,
  inTransit: 0,
  completed: 0,
  total: 0
})

const recentOrders = ref<any[]>([])

const loadOrders = async () => {
  try {
    const res: any = await getUserOrdersAPI()
    const orders = res.data || []
    recentOrders.value = orders.slice(0, 5)
    
    stats.value.total = orders.length
    stats.value.pending = orders.filter((o: any) => o.status === 0 || o.status === 1).length
    stats.value.inTransit = orders.filter((o: any) => o.status === 2 || o.status === 3).length
    stats.value.completed = orders.filter((o: any) => o.status === 4).length
  } catch (error) {
    console.error('加载订单失败:', error)
  }
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

const handleCreateOrder = () => {
  router.push('/user/order/create')
}

const handleViewSites = () => {
  router.push('/user/sites')
}

const handleViewOrders = () => {
  router.push('/user/orders')
}

const handleQueryLogistics = () => {
  router.push('/user/logistics')
}

onMounted(() => {
  loadOrders()
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

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.quick-actions .el-button {
  width: 100%;
  height: 50px;
}
</style>
