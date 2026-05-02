<template>
  <div class="warehouses-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>仓库管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加仓库
          </el-button>
        </div>
      </template>
      <el-table :data="warehouses" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="仓库名称" min-width="150" />
        <el-table-column prop="location" label="位置" min-width="200" />
        <el-table-column prop="capacity" label="容量" width="100" />
        <el-table-column prop="used" label="已用" width="100">
          <template #default="scope">
            <el-progress
              :percentage="scope.row.capacity > 0 ? Math.round((scope.row.used / scope.row.capacity) * 100) : 0"
              :stroke-width="20"
              :show-text="false"
            >
              <template #default="{ percentage }">
                <span class="percentage-value">{{ scope.row.used }}</span>
              </template>
            </el-progress>
          </template>
        </el-table-column>
        <el-table-column prop="manager" label="负责人" width="100" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '正常' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">
              编辑
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑仓库' : '添加仓库'" width="500px">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="仓库名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入仓库名称" />
        </el-form-item>
        <el-form-item label="位置" prop="location">
          <el-input v-model="formData.location" placeholder="请输入位置" />
        </el-form-item>
        <el-form-item label="容量" prop="capacity">
          <el-input-number v-model="formData.capacity" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="负责人" prop="manager">
          <el-input v-model="formData.manager" placeholder="请输入负责人" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="状态" prop="status" v-if="isEdit">
          <el-radio-group v-model="formData.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">关闭</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { getWarehousesAPI, createWarehouseAPI, updateWarehouseAPI } from '@/api'

const warehouses = ref<any[]>([])
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const loading = ref(false)
const isEdit = ref(false)
const editId = ref<number>(0)

const formData = reactive({
  name: '',
  location: '',
  capacity: 0,
  manager: '',
  phone: '',
  status: 1
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入仓库名称', trigger: 'blur' }]
}

const loadWarehouses = async () => {
  try {
    const res: any = await getWarehousesAPI()
    warehouses.value = res.data || []
  } catch (error) {
    console.error('加载仓库失败:', error)
  }
}

const resetForm = () => {
  formData.name = ''
  formData.location = ''
  formData.capacity = 0
  formData.manager = ''
  formData.phone = ''
  formData.status = 1
}

const handleAdd = () => {
  isEdit.value = false
  editId.value = 0
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  formData.name = row.name || ''
  formData.location = row.location || ''
  formData.capacity = row.capacity || 0
  formData.manager = row.manager || ''
  formData.phone = row.phone || ''
  formData.status = row.status || 1
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const data = {
          name: formData.name,
          location: formData.location,
          capacity: formData.capacity,
          manager: formData.manager,
          phone: formData.phone,
          status: formData.status
        }
        
        if (isEdit.value) {
          await updateWarehouseAPI(editId.value, data)
          ElMessage.success('更新成功')
        } else {
          await createWarehouseAPI(data)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        loadWarehouses()
      } catch (error) {
        console.error('提交失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadWarehouses()
})
</script>

<style scoped>
.warehouses-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.percentage-value {
  font-size: 12px;
}
</style>
