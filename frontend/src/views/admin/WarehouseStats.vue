<template>
  <div class="stats-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">仓库总数</span>
            <el-icon class="card-icon"><Box /></el-icon>
          </div>
          <div class="card-value">{{ stats.total_count || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">总容量</span>
            <el-icon class="card-icon"><Database /></el-icon>
          </div>
          <div class="card-value">{{ stats.total_capacity || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">已使用</span>
            <el-icon class="card-icon"><DocumentChecked /></el-icon>
          </div>
          <div class="card-value">{{ stats.total_used || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">使用率</span>
            <el-icon class="card-icon"><TrendCharts /></el-icon>
          </div>
          <div class="card-value">{{ stats.utilization ? stats.utilization.toFixed(2) : '0' }}%</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>仓库列表</span>
          </template>
          <el-table :data="warehouses" style="width: 100%">
            <el-table-column prop="name" label="仓库名称" min-width="150" />
            <el-table-column prop="location" label="位置" min-width="200" />
            <el-table-column prop="capacity" label="容量" width="120" />
            <el-table-column prop="used" label="已使用" width="120" />
            <el-table-column label="使用率" width="200">
              <template #default="scope">
                <el-progress
                  :percentage="scope.row.capacity > 0 ? Math.round((scope.row.used / scope.row.capacity) * 100) : 0"
                  :status="scope.row.capacity > 0 && (scope.row.used / scope.row.capacity) > 0.8 ? 'exception' : ''"
                />
              </template>
            </el-table-column>
            <el-table-column prop="manager" label="负责人" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                  {{ scope.row.status === 1 ? '正常' : '关闭' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { getWarehouseStatsAPI, getWarehousesAPI } from '@/api'

const stats = ref<any>({})
const warehouses = ref<any[]>([])

const loadData = async () => {
  try {
    const statsRes: any = await getWarehouseStatsAPI()
    stats.value = statsRes.data || {}
    
    const warehousesRes: any = await getWarehousesAPI()
    warehouses.value = warehousesRes.data || []
  } catch (error) {
    console.error('加载数据失败:', error)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.stats-container {
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
</style>
