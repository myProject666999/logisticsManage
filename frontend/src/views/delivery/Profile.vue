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
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="formData.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="身份证号" prop="idCard">
          <el-input v-model="formData.idCard" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="车牌号">
          <el-input v-model="formData.vehicleNo" placeholder="请输入车牌号" />
        </el-form-item>
        <el-form-item label="车辆类型">
          <el-select v-model="formData.vehicleType" placeholder="请选择车辆类型" style="width: 100%">
            <el-option label="电动车" :value="1" />
            <el-option label="摩托车" :value="2" />
            <el-option label="货车" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="审核状态">
          <el-tag :type="formData.verifyStatus === 1 ? 'success' : formData.verifyStatus === 2 ? 'danger' : 'warning'">
            {{ getVerifyStatusText(formData.verifyStatus) }}
          </el-tag>
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
  realName: '',
  idCard: '',
  vehicleNo: '',
  vehicleType: undefined as number | undefined,
  verifyStatus: 0
})

const rules: FormRules = {
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

const getVerifyStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待审核',
    1: '已通过',
    2: '已拒绝'
  }
  return map[status] || '未知'
}

const loadUserInfo = async () => {
  try {
    const res: any = await getUserInfoAPI('delivery')
    const data = res.data
    formData.username = data.username || ''
    formData.phone = data.phone || ''
    formData.realName = data.real_name || ''
    formData.idCard = data.id_card || ''
    formData.vehicleNo = data.vehicle_no || ''
    formData.vehicleType = data.vehicle_type
    formData.verifyStatus = data.verify_status || 0
  } catch (error) {
    console.error('加载配送员信息失败:', error)
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await updateUserInfoAPI('delivery', {
          phone: formData.phone,
          real_name: formData.realName,
          id_card: formData.idCard
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
