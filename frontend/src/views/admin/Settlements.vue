<template>
  <div class="settlements-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>结算管理</span>
          <el-select v-model="filterStatus" placeholder="全部状态" style="width: 150px" @change="loadSettlements">
            <el-option label="全部" value="" />
            <el-option label="待确认" value="0" />
            <el-option label="已确认" value="1" />
          </el-select>
        </div>
      </template>
      <el-table :data="settlements" style="width: 100%">
        <el-table-column prop="delivery_man_name" label="配送员" width="120" />
        <el-table-column prop="period_start" label="结算开始" width="180" />
        <el-table-column prop="period_end" label="结算结束" width="180" />
        <el-table-column prop="total_amount" label="总金额(元)" width="150">
          <template #default="scope">
            <span style="color: #f56c6c; font-weight: bold">¥{{ scope.row.total_amount }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="order_count" label="订单数" width="100" />
        <el-table-column prop="reimbursement_amount" label="报销(元)" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'warning'">
              {{ scope.row.status === 1 ? '已确认' : '待确认' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="confirm_time" label="确认时间" width="180">
          <template #default="scope">
            {{ scope.row.confirm_time || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <template v-if="scope.row.status === 0">
              <el-button type="success" link @click="handleConfirm(scope.row)">
                确认
              </el-button>
            </template>
            <template v-else>
              <span class="text-muted">已处理</span>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAllSettlementsAPI, confirmSettlementAPI } from '@/api'

const settlements = ref<any[]>([])
const filterStatus = ref('')

const loadSettlements = async () => {
  try {
    const params = filterStatus.value ? { status: filterStatus.value } : undefined
    const res: any = await getAllSettlementsAPI(params)
    settlements.value = res.data || []
  } catch (error) {
    console.error('加载结算失败:', error)
  }
}

const handleConfirm = (row: any) => {
  ElMessageBox.confirm(
    `确认结算配送员「${row.delivery_man_name}」的订单？\n总金额：¥${row.total_amount}`,
    '确认结算',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await confirmSettlementAPI(row.id)
      ElMessage.success('确认成功')
      loadSettlements()
    } catch (error) {
      console.error('确认失败:', error)
    }
  }).catch(() => {})
}

onMounted(() => {
  loadSettlements()
})
</script>

<style scoped>
.settlements-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.text-muted {
  color: #999;
}
</style>
