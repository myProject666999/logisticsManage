<template>
  <div class="delivery-men-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>配送员管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加配送员
          </el-button>
        </div>
      </template>
      <el-table :data="deliveryMen" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="real_name" label="真实姓名" width="120" />
        <el-table-column prop="id_card" label="身份证号" width="200" />
        <el-table-column prop="vehicle_no" label="车牌号" width="150" />
        <el-table-column prop="vehicle_type" label="车辆类型" width="120">
          <template #default="scope">
            {{ getVehicleTypeText(scope.row.vehicle_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="site_name" label="所属站点" width="150" />
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '在职' : '离职' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button 
              type="warning" 
              link 
              @click="handleToggleStatus(scope.row)"
            >
              {{ scope.row.status === 1 ? '禁用' : '启用' }}
            </el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑配送员' : '添加配送员'" width="600px">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="formData.username" placeholder="请输入用户名" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="formData.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="formData.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="身份证号" prop="idCard">
          <el-input v-model="formData.idCard" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="车牌号" prop="vehicleNo">
          <el-input v-model="formData.vehicleNo" placeholder="请输入车牌号" />
        </el-form-item>
        <el-form-item label="车辆类型" prop="vehicleType">
          <el-select v-model="formData.vehicleType" placeholder="请选择车辆类型" style="width: 100%">
            <el-option label="电动车" :value="1" />
            <el-option label="摩托车" :value="2" />
            <el-option label="货车" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属站点" prop="siteId">
          <el-select v-model="formData.siteId" placeholder="请选择所属站点" style="width: 100%">
            <el-option 
              v-for="site in sites" 
              :key="site.id" 
              :label="site.name" 
              :value="site.id" 
            />
          </el-select>
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
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { 
  getAllDeliveryMenAPI, 
  createDeliveryManAPI, 
  updateDeliveryManAPI, 
  deleteDeliveryManAPI,
  getAllSitesAPI
} from '@/api'

const deliveryMen = ref<any[]>([])
const sites = ref<any[]>([])
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const loading = ref(false)
const isEdit = ref(false)
const editId = ref<number>(0)

const formData = reactive({
  username: '',
  password: '',
  phone: '',
  realName: '',
  idCard: '',
  vehicleNo: '',
  vehicleType: undefined,
  siteId: undefined
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const getVehicleTypeText = (type: number) => {
  const map: Record<number, string> = {
    1: '电动车',
    2: '摩托车',
    3: '货车'
  }
  return map[type] || '其他'
}

const loadDeliveryMen = async () => {
  try {
    const res: any = await getAllDeliveryMenAPI()
    deliveryMen.value = res.data || []
  } catch (error) {
    console.error('加载配送员失败:', error)
  }
}

const loadSites = async () => {
  try {
    const res: any = await getAllSitesAPI()
    sites.value = res.data || []
  } catch (error) {
    console.error('加载站点失败:', error)
  }
}

const resetForm = () => {
  formData.username = ''
  formData.password = ''
  formData.phone = ''
  formData.realName = ''
  formData.idCard = ''
  formData.vehicleNo = ''
  formData.vehicleType = undefined
  formData.siteId = undefined
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
  formData.username = row.username || ''
  formData.phone = row.phone || ''
  formData.realName = row.real_name || ''
  formData.idCard = row.id_card || ''
  formData.vehicleNo = row.vehicle_no || ''
  formData.vehicleType = row.vehicle_type
  formData.siteId = row.site_id
  dialogVisible.value = true
}

const handleToggleStatus = async (row: any) => {
  const newStatus = row.status === 1 ? 0 : 1
  try {
    await updateDeliveryManAPI(row.id, { status: newStatus })
    ElMessage.success(newStatus === 1 ? '启用成功' : '禁用成功')
    loadDeliveryMen()
  } catch (error) {
    console.error('状态切换失败:', error)
  }
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除配送员 ${row.username} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteDeliveryManAPI(row.id)
    ElMessage.success('删除成功')
    loadDeliveryMen()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const data = {
          username: formData.username,
          password: formData.password,
          phone: formData.phone,
          real_name: formData.realName,
          id_card: formData.idCard,
          vehicle_no: formData.vehicleNo,
          vehicle_type: formData.vehicleType,
          site_id: formData.siteId
        }
        
        if (isEdit.value) {
          await updateDeliveryManAPI(editId.value, data)
          ElMessage.success('更新成功')
        } else {
          await createDeliveryManAPI(data)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        loadDeliveryMen()
      } catch (error) {
        console.error('提交失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadDeliveryMen()
  loadSites()
})
</script>

<style scoped>
.delivery-men-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
