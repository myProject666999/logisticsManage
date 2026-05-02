<template>
  <div class="orders-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的订单</span>
          <el-select v-model="filterStatus" placeholder="全部状态" style="width: 150px" @change="loadOrders">
            <el-option label="全部" value="" />
            <el-option label="待揽收" value="0" />
            <el-option label="已揽收" value="1" />
            <el-option label="运输中" value="2" />
            <el-option label="派送中" value="3" />
            <el-option label="已签收" value="4" />
          </el-select>
        </div>
      </template>
      <el-table :data="orders" style="width: 100%">
        <el-table-column prop="order_no" label="订单号" min-width="180" />
        <el-table-column prop="sender_name" label="寄件人" width="100" />
        <el-table-column prop="receiver_name" label="收件人" width="100" />
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
            <el-button type="primary" link @click="handleView(scope.row)">
              详情
            </el-button>
            <el-button v-if="scope.row.status >= 1 && scope.row.status < 4" type="warning" link @click="handleUpdateLogistics(scope.row)">
              更新物流
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="logisticsDialogVisible" title="更新物流信息" width="500px">
      <el-form
        ref="logisticsFormRef"
        :model="logisticsFormData"
        :rules="logisticsRules"
        label-width="100px"
      >
        <el-form-item label="订单号">
          <el-input :value="currentOrder?.order_no" disabled />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="logisticsFormData.status" placeholder="请选择状态" style="width: 100%">
            <el-option label="已揽收" :value="1" />
            <el-option label="运输中" :value="2" />
            <el-option label="派送中" :value="3" />
            <el-option label="已签收" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="logisticsFormData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入物流描述"
          />
        </el-form-item>
        <el-form-item label="地点" prop="location">
          <el-input v-model="logisticsFormData.location" placeholder="请输入当前地点" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="logisticsDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitLogistics" :loading="logisticsLoading">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { getDeliveryOrdersAPI, updateLogisticsAPI } from '@/api'

const orders = ref<any[]>([])
const filterStatus = ref('')
const logisticsDialogVisible = ref(false)
const logisticsFormRef = ref<FormInstance>()
const logisticsLoading = ref(false)
const currentOrder = ref<any>(null)

const logisticsFormData = reactive({
  status: 1,
  description: '',
  location: ''
})

const logisticsRules: FormRules = {
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
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

const loadOrders = async () => {
  try {
    const params = filterStatus.value ? { status: filterStatus.value } : undefined
    const res: any = await getDeliveryOrdersAPI(params)
    orders.value = res.data || []
  } catch (error) {
    console.error('加载订单失败:', error)
  }
}

const handleView = (row: any) => {
  ElMessage.info('订单号: ' + row.order_no)
}

const handleUpdateLogistics = (row: any) => {
  currentOrder.value = row
  logisticsFormData.status = row.status > 1 ? row.status : 2
  logisticsFormData.description = ''
  logisticsFormData.location = ''
  logisticsDialogVisible.value = true
}

const submitLogistics = async () => {
  if (!logisticsFormRef.value || !currentOrder.value) return
  
  await logisticsFormRef.value.validate(async (valid) => {
    if (valid) {
      logisticsLoading.value = true
      try {
        await updateLogisticsAPI(currentOrder.value.order_no, {
          status: logisticsFormData.status,
          description: logisticsFormData.description,
          location: logisticsFormData.location
        })
        ElMessage.success('物流信息更新成功')
        logisticsDialogVisible.value = false
        loadOrders()
      } catch (error) {
        console.error('更新失败:', error)
      } finally {
        logisticsLoading.value = false
      }
    }
  })
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
