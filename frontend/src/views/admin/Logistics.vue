<template>
  <div class="logistics-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>物流管理</span>
          <el-input
            v-model="searchOrderNo"
            placeholder="请输入订单号搜索"
            style="width: 250px"
            @keyup.enter="loadLogistics"
          >
            <template #append>
              <el-button @click="loadLogistics">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </div>
      </template>
      <el-table :data="logisticsList" style="width: 100%">
        <el-table-column prop="order_no" label="订单号" min-width="180" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="location" label="地点" min-width="150" />
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getLogisticsListAPI } from '@/api'

const logisticsList = ref<any[]>([])
const searchOrderNo = ref('')

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

const loadLogistics = async () => {
  try {
    const params = searchOrderNo.value ? { order_no: searchOrderNo.value } : undefined
    const res: any = await getLogisticsListAPI(params)
    logisticsList.value = res.data || []
  } catch (error) {
    console.error('加载物流信息失败:', error)
  }
}

onMounted(() => {
  loadLogistics()
})
</script>

<style scoped>
.logistics-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
