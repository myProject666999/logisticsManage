<template>
  <div class="orders-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>订单管理</span>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索订单号/寄件人/收件人"
            style="width: 250px"
            clearable
            @clear="loadOrders"
            @keyup.enter="loadOrders"
          >
            <template #append>
              <el-button @click="loadOrders">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
          <el-select v-model="filterStatus" placeholder="全部状态" style="width: 120px" @change="loadOrders" clearable>
            <el-option label="待揽收" :value="0" />
            <el-option label="已揽收" :value="1" />
            <el-option label="运输中" :value="2" />
            <el-option label="派送中" :value="3" />
            <el-option label="已签收" :value="4" />
            <el-option label="已取消" :value="5" />
          </el-select>
        </div>
      </template>
      <el-table :data="orders" style="width: 100%">
        <el-table-column prop="order_no" label="订单号" min-width="180" />
        <el-table-column prop="sender_name" label="寄件人" width="100" />
        <el-table-column prop="sender_phone" label="寄件电话" width="130" />
        <el-table-column prop="receiver_name" label="收件人" width="100" />
        <el-table-column prop="receiver_addr" label="收件地址" min-width="200" />
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
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="handleView(scope.row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getAdminOrdersAPI } from '@/api'

const orders = ref<any[]>([])
const searchKeyword = ref('')
const filterStatus = ref<number | ''>('')

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

const loadOrders = async () => {
  try {
    const params: any = {}
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    if (filterStatus.value !== '') {
      params.status = String(filterStatus.value)
    }
    const res: any = await getAdminOrdersAPI(params)
    orders.value = res.data || []
  } catch (error) {
    console.error('加载订单失败:', error)
  }
}

const handleView = (row: any) => {
  ElMessage.info('订单号: ' + row.order_no)
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.orders-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}
</style>
