<template>
  <div class="stats-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">报销总数</span>
            <el-icon class="card-icon"><Document /></el-icon>
          </div>
          <div class="card-value">{{ stats.total || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">待审核</span>
            <el-icon class="card-icon"><Clock /></el-icon>
          </div>
          <div class="card-value">{{ stats.pending || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="card-item">
          <div class="card-header">
            <span class="card-title">已通过</span>
            <el-icon class="card-icon"><CircleCheck /></el-icon>
          </div>
          <div class="card-value">{{ stats.approved || 0 }}</div>
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
            <span>报销类型分布</span>
          </template>
          <el-table :data="typeStats" style="width: 100%">
            <el-table-column prop="type" label="类型" width="120">
              <template #default="scope">
                {{ getTypeText(scope.row.type) }}
              </template>
            </el-table-column>
            <el-table-column prop="count" label="数量" width="100" />
            <el-table-column prop="amount" label="金额(元)" width="150" />
            <el-table-column label="占比" min-width="200">
              <template #default="scope">
                <el-progress 
                  :percentage="stats.total > 0 ? Math.round((scope.row.count / stats.total) * 100) : 0"
                />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>配送员报销排行</span>
          </template>
          <el-table :data="deliverymanStats" style="width: 100%">
            <el-table-column prop="delivery_man_name" label="配送员" width="120" />
            <el-table-column prop="count" label="报销次数" width="100" />
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
import { getReimbursementStatsAPI } from '@/api'

const stats = ref<any>({})
const typeStats = ref<any[]>([])
const deliverymanStats = ref<any[]>([])

const getTypeText = (type: number) => {
  const map: Record<number, string> = {
    1: '油费',
    2: '过路费',
    3: '维修费',
    4: '其他'
  }
  return map[type] || '其他'
}

const loadData = async () => {
  try {
    const res: any = await getReimbursementStatsAPI()
    stats.value = res.data || {}
    typeStats.value = res.data?.by_type || []
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
</style>
