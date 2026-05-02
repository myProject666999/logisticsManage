<template>
  <div class="pending-orders-container">
    <el-card>
      <template #header>
        <span>待揽收订单</span>
      </template>
      <el-table :data="orders" style="width: 100%">
        <el-table-column prop="order_no" label="订单号" min-width="180" />
        <el-table-column prop="sender_name" label="寄件人" width="100" />
        <el-table-column prop="sender_phone" label="寄件电话" width="130" />
        <el-table-column prop="sender_address" label="寄件地址" min-width="200" />
        <el-table-column prop="receiver_name" label="收件人" width="100" />
        <el-table-column prop="receiver_addr" label="收件地址" min-width="200" />
        <el-table-column prop="total_fee" label="运费(元)" width="100">
          <template #default="scope">
            ¥{{ scope.row.total_fee }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" @click="handlePickup(scope.row)">
              揽收
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
import { getPendingOrdersAPI, pickupOrderAPI } from '@/api'

const orders = ref<any[]>([])

const loadOrders = async () => {
  try {
    const res: any = await getPendingOrdersAPI()
    orders.value = res.data || []
  } catch (error) {
    console.error('加载订单失败:', error)
  }
}

const handlePickup = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要揽收订单 ${row.order_no} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await pickupOrderAPI(row.order_no)
    ElMessage.success('揽收成功')
    loadOrders()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('揽收失败:', error)
    }
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.pending-orders-container {
  padding: 0;
}
</style>
