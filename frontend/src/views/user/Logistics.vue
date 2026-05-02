<template>
  <div class="logistics-container">
    <el-card>
      <template #header>
        <span>物流查询</span>
      </template>
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="订单号">
          <el-input
            v-model="searchForm.orderNo"
            placeholder="请输入订单号"
            style="width: 300px"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch" :loading="loading">
            查询
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card v-if="order" style="margin-top: 20px">
      <template #header>
        <span>订单信息</span>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(order.status)">
            {{ getStatusText(order.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="寄件人">{{ order.sender_name }}</el-descriptions-item>
        <el-descriptions-item label="收件人">{{ order.receiver_name }}</el-descriptions-item>
        <el-descriptions-item label="寄件地址">{{ order.sender_address }}</el-descriptions-item>
        <el-descriptions-item label="收件地址">{{ order.receiver_addr }}</el-descriptions-item>
        <el-descriptions-item label="物品类型">{{ order.goods_type || '普通物品' }}</el-descriptions-item>
        <el-descriptions-item label="运费">¥{{ order.total_fee }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card v-if="logistics.length > 0" style="margin-top: 20px">
      <template #header>
        <span>物流跟踪</span>
      </template>
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
          <p v-if="item.operator">操作人: {{ item.operator }}</p>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <el-empty v-else-if="searched && !order" description="未找到该订单的物流信息" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { queryLogisticsAPI } from '@/api'

const loading = ref(false)
const searched = ref(false)
const order = ref<any>(null)
const logistics = ref<any[]>([])

const searchForm = reactive({
  orderNo: ''
})

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

const handleSearch = async () => {
  if (!searchForm.orderNo.trim()) {
    ElMessage.warning('请输入订单号')
    return
  }
  
  loading.value = true
  searched.value = true
  order.value = null
  logistics.value = []
  
  try {
    const res: any = await queryLogisticsAPI(searchForm.orderNo)
    order.value = res.data.order
    logistics.value = res.data.logistics || []
  } catch (error) {
    console.error('查询物流失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.logistics-container {
  padding: 0;
}
</style>
