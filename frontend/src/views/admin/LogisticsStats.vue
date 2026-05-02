<template>
  <div class="stats-container">
    <el-row :gutter="20">
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">订单总数</span>
            <el-icon class="card-icon"><Document /></el-icon>
          </div>
          <div class="card-value">{{ stats.total || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">待揽收</span>
            <el-icon class="card-icon"><Clock /></el-icon>
          </div>
          <div class="card-value">{{ stats.pending || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">运输中</span>
            <el-icon class="card-icon"><Truck /></el-icon>
          </div>
          <div class="card-value">{{ stats.in_transit || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">已签收</span>
            <el-icon class="card-icon"><CircleCheck /></el-icon>
          </div>
          <div class="card-value">{{ stats.completed || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">已取消</span>
            <el-icon class="card-icon"><CircleClose /></el-icon>
          </div>
          <div class="card-value">{{ stats.cancelled || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">完成率</span>
            <el-icon class="card-icon"><TrendCharts /></el-icon>
          </div>
          <div class="card-value">
            {{ stats.total > 0 ? ((stats.completed / stats.total) * 100).toFixed(1) : 0 }}%
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>订单状态分布</span>
          </template>
          <div class="chart-placeholder">
            <el-row :gutter="40">
              <el-col :span="6">
                <div class="stat-item">
                  <div class="stat-label">待揽收</div>
                  <el-progress 
                    :percentage="stats.total > 0 ? Math.round((stats.pending / stats.total) * 100) : 0"
                    status="warning"
                  />
                </div>
              </el-col>
              <el-col :span="6">
                <div class="stat-item">
                  <div class="stat-label">运输中</div>
                  <el-progress 
                    :percentage="stats.total > 0 ? Math.round((stats.in_transit / stats.total) * 100) : 0"
                    status="exception"
                  />
                </div>
              </el-col>
              <el-col :span="6">
                <div class="stat-item">
                  <div class="stat-label">已签收</div>
                  <el-progress 
                    :percentage="stats.total > 0 ? Math.round((stats.completed / stats.total) * 100) : 0"
                    status="success"
                  />
                </div>
              </el-col>
              <el-col :span="6">
                <div class="stat-item">
                  <div class="stat-label">已取消</div>
                  <el-progress 
                    :percentage="stats.total > 0 ? Math.round((stats.cancelled / stats.total) * 100) : 0"
                  />
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { getLogisticsStatsAPI } from '@/api'

const stats = ref<any>({})

const loadData = async () => {
  try {
    const res: any = await getLogisticsStatsAPI()
    stats.value = res.data || {}
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

.stat-item {
  text-align: center;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}
</style>
