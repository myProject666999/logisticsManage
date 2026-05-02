<template>
  <div class="stats-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">结算总数</span>
            <el-icon class="card-icon"><Document /></el-icon>
          </div>
          <div class="card-value">{{ stats.total || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">待确认</span>
            <el-icon class="card-icon"><Clock /></el-icon>
          </div>
          <div class="card-value">{{ stats.pending || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">已确认</span>
            <el-icon class="card-icon"><CircleCheck /></el-icon>
          </div>
          <div class="card-value">{{ stats.confirmed || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">总金额</span>
            <el-icon class="card-icon"><Money /></el-icon>
          </div>
          <div class="card-value">¥{{ stats.total_amount || 0 }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>本月结算趋势</span>
          </template>
          <div class="chart-placeholder">
            <el-table :data="monthlyStats" style="width: 100%">
              <el-table-column prop="month" label="月份" width="100" />
              <el-table-column prop="count" label="结算数量" width="120" />
              <el-table-column prop="amount" label="结算金额" width="150">
                <template #default="scope">
                  ¥{{ scope.row.amount }}
                </template>
              </el-table-column>
              <el-table-column label="金额趋势" min-width="200">
                <template #default="scope">
                  <el-progress 
                    :percentage="stats.total_amount > 0 ? Math.round((scope.row.amount / stats.total_amount) * 100) : 0"
                    status="success"
                  />
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>配送员结算排行</span>
          </template>
          <el-table :data="deliverymanStats" style="width: 100%">
            <el-table-column prop="delivery_man_name" label="配送员" width="120" />
            <el-table-column prop="count" label="结算次数" width="100" />
            <el-table-column prop="amount" label="总金额(元)" width="150">
              <template #default="scope">
                <span style="color: #f56c6c; font-weight: bold">¥{{ scope.row.amount }}</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSettlementStatsAPI } from '@/api'

const stats = ref<any>({})
const monthlyStats = ref<any[]>([])
const deliverymanStats = ref<any[]>([])

const loadData = async () => {
  try {
    const res: any = await getSettlementStatsAPI()
    stats.value = res.data || {}
    monthlyStats.value = res.data?.monthly || []
    deliverymanStats.value = res.data?.by_deliveryman || []
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

.chart-placeholder {
  padding: 20px;
}
</style>
