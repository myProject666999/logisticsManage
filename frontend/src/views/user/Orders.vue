<template>
  <div class="orders-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>订单列表</span>
          <el-select v-model="statusFilter" placeholder="全部状态" clearable @change="loadOrders">
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
        <el-table-column prop="sender_name" label="寄件人" min-width="100" />
        <el-table-column prop="sender_address" label="寄件地址" min-width="200" />
        <el-table-column prop="receiver_name" label="收件人" min-width="100" />
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
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="viewDetail(scope.row)">
              详情
            </el-button>
            <el-button 
              type="danger" 
              link 
              :disabled="scope.row.status !== 0"
              @click="handleCancel(scope.row)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="detailVisible" title="订单详情" width="600px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="订单号">{{ currentOrder?.order_no }}</el-descriptions-item>
        <el-descriptions-item label="寄件人">{{ currentOrder?.sender_name }}</el-descriptions-item>
        <el-descriptions-item label="寄件电话">{{ currentOrder?.sender_phone }}</el-descriptions-item>
        <el-descriptions-item label="寄件地址">{{ currentOrder?.sender_address }}</el-descriptions-item>
        <el-descriptions-item label="收件人">{{ currentOrder?.receiver_name }}</el-descriptions-item>
        <el-descriptions-item label="收件电话">{{ currentOrder?.receiver_phone }}</el-descriptions-item>
        <el-descriptions-item label="收件地址">{{ currentOrder?.receiver_addr }}</el-descriptions-item>
        <el-descriptions-item label="物品类型">{{ currentOrder?.goods_type || '普通物品' }}</el-descriptions-item>
        <el-descriptions-item label="重量">{{ currentOrder?.goods_weight }}kg</el-descriptions-item>
        <el-descriptions-item label="数量">{{ currentOrder?.goods_amount }}件</el-descriptions-item>
        <el-descriptions-item label="运费">¥{{ currentOrder?.total_fee }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentOrder?.status || 0)">
            {{ getStatusText(currentOrder?.status || 0) }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
      
      <div v-if="logistics.length > 0" style="margin-top: 20px">
        <h4 style="margin-bottom: 15px">物流信息</h4>
        <el-timeline>
          <el-timeline-item
            v-for="(item, index) in logistics"
            :key="item.id"
            :timestamp="item.create_time"
            placement="top"
            :type="index === 0 ? 'primary' : ''"
          >
            <h4>{{ item.description }}</h4>
            <p v-if="item.location">位置: {{ item.location }}</p>
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserOrdersAPI, getUserOrderDetailAPI, cancelOrderAPI } from '@/api'

const statusFilter = ref<string | undefined>(undefined)
const orders = ref<any[]>([])
const detailVisible = ref(false)
const currentOrder = ref<any>(null)
const logistics = ref<any[]>([])

const loadOrders = async () => {
  try {
    const params = statusFilter.value !== undefined ? { status: String(statusFilter.value) } : undefined
    const res: any = await getUserOrdersAPI(params)
    orders.value = res.data || []
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

const viewDetail = async (row: any) => {
  try {
    const res: any = await getUserOrderDetailAPI(row.order_no)
    currentOrder.value = res.data.order
    logistics.value = res.data.logistics || []
    detailVisible.value = true
  } catch (error) {
    console.error('加载订单详情失败:', error)
  }
}

const handleCancel = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await cancelOrderAPI(row.order_no)
    ElMessage.success('订单已取消')
    loadOrders()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('取消订单失败:', error)
    }
  }
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
}
</style>
