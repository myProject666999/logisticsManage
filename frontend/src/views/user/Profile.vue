<template>
  <div class="profile-container">
    <el-card>
      <template #header>
        <span>个人信息</span>
      </template>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
        style="max-width: 500px"
      >
        <el-form-item label="用户名">
          <el-input v-model="formData.username" disabled />
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
        <el-form-item label="地址">
          <el-input 
            v-model="formData.address" 
            type="textarea"
            :rows="3"
            placeholder="请输入地址"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            保存修改
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { getUserInfoAPI, updateUserInfoAPI } from '@/api'

const formRef = ref<FormInstance>()
const loading = ref(false)

const formData = reactive({
  username: '',
  phone: '',
  email: '',
  realName: '',
  address: ''
})

const rules: FormRules = {
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const loadUserInfo = async () => {
  try {
    const res: any = await getUserInfoAPI('user')
    const data = res.data
    formData.username = data.username || ''
    formData.phone = data.phone || ''
    formData.email = data.email || ''
    formData.realName = data.real_name || ''
    formData.address = data.address || ''
  } catch (error) {
    console.error('加载用户信息失败:', error)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await updateUserInfoAPI('user', {
          phone: formData.phone,
          email: formData.email,
          real_name: formData.realName,
          address: formData.address
        })
        ElMessage.success('修改成功')
      } catch (error) {
        console.error('修改失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadUserInfo()
})
</script>

<style scoped>
.profile-container {
  padding: 0;
}
</style>
