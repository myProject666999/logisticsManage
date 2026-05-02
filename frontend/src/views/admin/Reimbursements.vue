<template>
  <div class="reimbursements-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>报销管理</span>
          <el-select v-model="filterStatus" placeholder="全部状态" style="width: 150px" @change="loadReimbursements">
            <el-option label="全部" value="" />
            <el-option label="待审核" value="0" />
            <el-option label="已通过" value="1" />
            <el-option label="已拒绝" value="2" />
          </el-select>
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
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <template v-if="scope.row.status === 0">
              <el-button type="success" link @click="handleAudit(scope.row, 1)">
                通过
              </el-button>
              <el-button type="danger" link @click="handleAudit(scope.row, 2)">
                拒绝
              </el-button>
            </template>
            <template v-else>
              <span class="text-muted">已审核</span>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="auditDialogVisible" title="审核报销" width="400px">
      <el-form
        ref="auditFormRef"
        :model="auditFormData"
        :rules="auditRules"
        label-width="80px"
      >
        <el-form-item label="订单号">
          <el-input :value="currentReimbursement?.order_no || '无'" disabled />
        </el-form-item>
        <el-form-item label="类型">
          <el-input :value="getTypeText(currentReimbursement?.type)" disabled />
        </el-form-item>
        <el-form-item label="金额">
          <el-input :value="'¥' + (currentReimbursement?.amount || 0)" disabled />
        </el-form-item>
        <el-form-item label="审核结果" prop="status">
          <el-radio-group v-model="auditFormData.status">
            <el-radio :value="1">通过</el-radio>
            <el-radio :value="2">拒绝</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="auditFormData.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入审核备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="auditDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAudit" :loading="auditLoading">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { getAllReimbursementsAPI, auditReimbursementAPI } from '@/api'

const reimbursements = ref<any[]>([])
const filterStatus = ref('')
const auditDialogVisible = ref(false)
const auditFormRef = ref<FormInstance>()
const auditLoading = ref(false)
const currentReimbursement = ref<any>(null)

const auditFormData = reactive({
  status: 1,
  remark: ''
})

const auditRules: FormRules = {
  status: [{ required: true, message: '请选择审核结果', trigger: 'change' }]
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
    const params = filterStatus.value ? { status: filterStatus.value } : undefined
    const res: any = await getAllReimbursementsAPI(params)
    reimbursements.value = res.data || []
  } catch (error) {
    console.error('加载报销失败:', error)
  }
}

const handleAudit = (row: any, status: number) => {
  currentReimbursement.value = row
  auditFormData.status = status
  auditFormData.remark = ''
  auditDialogVisible.value = true
}

const submitAudit = async () => {
  if (!auditFormRef.value || !currentReimbursement.value) return
  
  await auditFormRef.value.validate(async (valid) => {
    if (valid) {
      auditLoading.value = true
      try {
        await auditReimbursementAPI(currentReimbursement.value.id, {
          status: auditFormData.status,
          remark: auditFormData.remark
        })
        ElMessage.success('审核成功')
        auditDialogVisible.value = false
        loadReimbursements()
      } catch (error) {
        console.error('审核失败:', error)
      } finally {
        auditLoading.value = false
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

.text-muted {
  color: #999;
}
</style>
