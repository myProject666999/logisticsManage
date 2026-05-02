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
        label-width="120px"
        style="max-width: 500px"
      >
        <el-form-item label="用户名">
          <el-input v-model="formData.username" disabled />
        </el-form-item>
        <el-form-item label="姓名" prop="realName">
          <el-input v-model="formData.realName" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="formData.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="formData.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="身份证号">
          <el-input v-model="formData.idCard" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSave" :loading="loading">
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
import { getAdminProfileAPI, updateAdminProfileAPI } from '@/api'
import { useUserStore } from '@/store'

const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const formData = reactive({
  username: '',
  realName: '',
  phone: '',
  email: '',
  idCard: ''
})

const rules: FormRules = {
  realName: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const loadProfile = async () => {
  try {
    const res: any = await getAdminProfileAPI()
    const data = res.data || {}
    formData.username = data.username || userStore.username || ''
    formData.realName = data.real_name || ''
    formData.phone = data.phone || ''
    formData.email = data.email || ''
    formData.idCard = data.id_card || ''
  } catch (error) {
    console.error('加载个人信息失败:', error)
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const data = {
          real_name: formData.realName,
          phone: formData.phone,
          email: formData.email,
          id_card: formData.idCard
        }
        await updateAdminProfileAPI(data)
        ElMessage.success('保存成功')
      } catch (error) {
        console.error('保存失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-container {
  padding: 0;
}
</style>
