<template>
  <div class="reimbursements-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>报销管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            申请报销
          </el-button>
        </div>
      </template>
      <el-table :data="reimbursements" style="width: 100%">
        <el-table-column prop="order_no" label="关联订单" width="180" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="scope">
            {{ getTypeText(scope.row.type) }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额(元)" width="120">
          <template #default="scope">
            ¥{{ scope.row.amount }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="说明" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="apply_time" label="申请时间" width="180" />
        <el-table-column prop="admin_remark" label="审核备注" min-width="150" />
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="申请报销" width="500px">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="关联订单" prop="orderNo">
          <el-input v-model="formData.orderNo" placeholder="请输入订单号（可选）" />
        </el-form-item>
        <el-form-item label="报销类型" prop="type">
          <el-select v-model="formData.type" placeholder="请选择报销类型" style="width: 100%">
            <el-option label="油费" :value="1" />
            <el-option label="过路费" :value="2" />
            <el-option label="维修费" :value="3" />
            <el-option label="其他" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input-number v-model="formData.amount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="说明" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入报销说明"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { getReimbursementsAPI, applyReimbursementAPI } from '@/api'

const reimbursements = ref<any[]>([])
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const loading = ref(false)

const formData = reactive({
  orderNo: '',
  type: undefined,
  amount: 0,
  description: ''
})

const rules: FormRules = {
  type: [{ required: true, message: '请选择报销类型', trigger: 'change' }],
  amount: [{ required: true, message: '请输入金额', trigger: 'blur' }]
}

const getTypeText = (type: number) => {
  const map: Record<number, string> = {
    1: '油费',
    2: '过路费',
    3: '维修费',
    4: '其他'
  }
  return map[type] || '其他'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待审核',
    1: '已通过',
    2: '已拒绝'
  }
  return map[status] || '未知'
}

const getStatusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'warning',
    1: 'success',
    2: 'danger'
  }
  return map[status] || 'info'
}

const loadReimbursements = async () => {
  try {
    const res: any = await getReimbursementsAPI()
    reimbursements.value = res.data || []
  } catch (error) {
    console.error('加载报销失败:', error)
  }
}

const handleAdd = () => {
  formData.orderNo = ''
  formData.type = undefined
  formData.amount = 0
  formData.description = ''
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await applyReimbursementAPI({
          order_no: formData.orderNo,
          type: formData.type,
          amount: formData.amount,
          description: formData.description
        })
        ElMessage.success('申请提交成功')
        dialogVisible.value = false
        loadReimbursements()
      } catch (error) {
        console.error('提交失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadReimbursements()
})
</script>

<style scoped>
.reimbursements-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
