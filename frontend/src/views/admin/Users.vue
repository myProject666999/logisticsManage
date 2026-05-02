<template>
  <div class="users-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加用户
          </el-button>
        </div>
      </template>
      <el-table :data="users" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="real_name" label="真实姓名" width="120" />
        <el-table-column prop="id_card" label="身份证号" width="200" />
        <el-table-column prop="address" label="地址" min-width="200" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button type="primary" link @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button type="danger" link @click="handleDelete(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑用户' : '添加用户'" width="600px">
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
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="formData.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="身份证号" prop="idCard">
          <el-input v-model="formData.idCard" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="formData.address" placeholder="请输入地址" />
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
import { getAllUsersAPI, createUserAPI, updateUserAPI, deleteUserAPI } from '@/api'

const users = ref<any[]>([])
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const loading = ref(false)
const isEdit = ref(false)
const editId = ref<number>(0)

const formData = reactive({
  username: '',
  password: '',
  phone: '',
  email: '',
  realName: '',
  idCard: '',
  address: ''
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const loadUsers = async () => {
  try {
    const res: any = await getAllUsersAPI()
    users.value = res.data || []
  } catch (error) {
    console.error('加载用户失败:', error)
  }
}

const resetForm = () => {
  formData.username = ''
  formData.password = ''
  formData.phone = ''
  formData.email = ''
  formData.realName = ''
  formData.idCard = ''
  formData.address = ''
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
  formData.email = row.email || ''
  formData.realName = row.real_name || ''
  formData.idCard = row.id_card || ''
  formData.address = row.address || ''
  dialogVisible.value = true
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除用户 ${row.username} 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteUserAPI(row.id)
    ElMessage.success('删除成功')
    loadUsers()
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
          email: formData.email,
          real_name: formData.realName,
          id_card: formData.idCard,
          address: formData.address
        }
        
        if (isEdit.value) {
          await updateUserAPI(editId.value, data)
          ElMessage.success('更新成功')
        } else {
          await createUserAPI(data)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        loadUsers()
      } catch (error) {
        console.error('提交失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.users-container {
  padding: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
