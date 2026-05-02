<template>
  <div class="sites-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>站点管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加站点
          </el-button>
        </div>
      </template>
      <el-table :data="sites" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="站点名称" min-width="150" />
        <el-table-column prop="province" label="省份" width="100" />
        <el-table-column prop="city" label="城市" width="100" />
        <el-table-column prop="district" label="区县" width="100" />
        <el-table-column prop="address" label="地址" min-width="200" />
        <el-table-column prop="phone" label="联系电话" width="130" />
        <el-table-column prop="manager" label="负责人" width="100" />
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

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑站点' : '添加站点'" width="600px">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="站点名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入站点名称" />
        </el-form-item>
        <el-form-item label="省份" prop="province">
          <el-input v-model="formData.province" placeholder="请输入省份" />
        </el-form-item>
        <el-form-item label="城市" prop="city">
          <el-input v-model="formData.city" placeholder="请输入城市" />
        </el-form-item>
        <el-form-item label="区县" prop="district">
          <el-input v-model="formData.district" placeholder="请输入区县" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="formData.address" placeholder="请输入详细地址" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="负责人" prop="manager">
          <el-input v-model="formData.manager" placeholder="请输入负责人" />
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
import { getAdminSitesAPI, createSiteAPI, updateSiteAPI } from '@/api'

const sites = ref<any[]>([])
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const loading = ref(false)
const isEdit = ref(false)
const editId = ref<number>(0)

const formData = reactive({
  name: '',
  province: '',
  city: '',
  district: '',
  address: '',
  phone: '',
  manager: '',
  status: 1
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入站点名称', trigger: 'blur' }],
  address: [{ required: true, message: '请输入地址', trigger: 'blur' }]
}

const loadSites = async () => {
  try {
    const res: any = await getAdminSitesAPI()
    sites.value = res.data || []
  } catch (error) {
    console.error('加载站点失败:', error)
  }
}

const resetForm = () => {
  formData.name = ''
  formData.province = ''
  formData.city = ''
  formData.district = ''
  formData.address = ''
  formData.phone = ''
  formData.manager = ''
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
  formData.province = row.province || ''
  formData.city = row.city || ''
  formData.district = row.district || ''
  formData.address = row.address || ''
  formData.phone = row.phone || ''
  formData.manager = row.manager || ''
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
          province: formData.province,
          city: formData.city,
          district: formData.district,
          address: formData.address,
          phone: formData.phone,
          manager: formData.manager,
          status: formData.status
        }
        
        if (isEdit.value) {
          await updateSiteAPI(editId.value, data)
          ElMessage.success('更新成功')
        } else {
          await createSiteAPI(data)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        loadSites()
      } catch (error) {
        console.error('提交失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadSites()
})
</script>

<style scoped>
.sites-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
