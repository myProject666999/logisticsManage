<template>
  <div class="delivery-info-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>配送信息管理</span>
          <el-select v-model="filterStatus" placeholder="全部状态" style="width: 150px" @change="loadDeliveries">
            <el-option label="全部" value="" />
            <el-option label="待配送" value="0" />
            <el-option label="配送中" value="1" />
            <el-option label="已完成" value="2" />
            <el-option label="已退回" value="3" />
          </el-select>
        </div>
      </template>
      <el-table :data="deliveries" style="width: 100%">
        <el-table-column prop="order_no" label="订单号" min-width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="start_time" label="开始时间" width="180" />
        <el-table-column prop="end_time" label="结束时间" width="180" />
        <el-table-column prop="notes" label="备注" min-width="200" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button 
              v-if="scope.row.status === 0" 
              type="primary" 
              link 
              @click="handleStartDelivery(scope.row)"
            >
              开始配送
            </el-button>
            <el-button 
              v-else-if="scope.row.status === 1" 
              type="success" 
              link 
              @click="handleCompleteDelivery(scope.row)"
            >
              完成配送
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDeliveriesAPI, updateDeliveryStatusAPI } from '@/api'

const deliveries = ref<any[]>([])
const filterStatus = ref('')

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待配送',
    1: '配送中',
    2: '已完成',
    3: '已退回'
  }
  return map[status] || '未知'
}

const getStatusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'warning',
    1: 'primary',
    2: 'success',
    3: 'danger'
  }
  return map[status] || 'info'
}

const loadDeliveries = async () => {
  try {
    const params = filterStatus.value ? { status: filterStatus.value } : undefined
    const res: any = await getDeliveriesAPI(params)
    deliveries.value = res.data || []
  } catch (error) {
    console.error('加载配送信息失败:', error)
  }
}

const handleStartDelivery = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要开始配送吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await updateDeliveryStatusAPI(row.order_no, { status: 1 })
    ElMessage.success('已开始配送')
    loadDeliveries()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  }
}

const handleCompleteDelivery = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要完成配送吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'success'
    })
    
    await updateDeliveryStatusAPI(row.order_no, { status: 2 })
    ElMessage.success('配送完成')
    loadDeliveries()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('操作失败:', error)
    }
  }
}

onMounted(() => {
  loadDeliveries()
})
</script>

<style scoped>
.delivery-info-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
