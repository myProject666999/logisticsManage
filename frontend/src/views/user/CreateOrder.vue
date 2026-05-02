<template>
  <div class="create-order-container">
    <el-card>
      <template #header>
        <span>寄件</span>
      </template>
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
        style="max-width: 600px"
      >
        <el-divider content-position="left">寄件人信息</el-divider>
        <el-form-item label="姓名" prop="senderName">
          <el-input v-model="formData.senderName" placeholder="请输入寄件人姓名" />
        </el-form-item>
        <el-form-item label="电话" prop="senderPhone">
          <el-input v-model="formData.senderPhone" placeholder="请输入寄件人电话" />
        </el-form-item>
        <el-form-item label="地址" prop="senderAddress">
          <el-input
            v-model="formData.senderAddress"
            type="textarea"
            :rows="2"
            placeholder="请输入寄件人地址"
          />
        </el-form-item>
        
        <el-divider content-position="left">收件人信息</el-divider>
        <el-form-item label="姓名" prop="receiverName">
          <el-input v-model="formData.receiverName" placeholder="请输入收件人姓名" />
        </el-form-item>
        <el-form-item label="电话" prop="receiverPhone">
          <el-input v-model="formData.receiverPhone" placeholder="请输入收件人电话" />
        </el-form-item>
        <el-form-item label="地址" prop="receiverAddr">
          <el-input
            v-model="formData.receiverAddr"
            type="textarea"
            :rows="2"
            placeholder="请输入收件人地址"
          />
        </el-form-item>
        
        <el-divider content-position="left">物品信息</el-divider>
        <el-form-item label="物品类型" prop="goodsType">
          <el-select v-model="formData.goodsType" placeholder="请选择物品类型" style="width: 200px">
            <el-option label="文件" value="文件" />
            <el-option label="电子产品" value="电子产品" />
            <el-option label="服装" value="服装" />
            <el-option label="食品" value="食品" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="重量(kg)" prop="goodsWeight">
          <el-input-number v-model="formData.goodsWeight" :min="0" :precision="2" :step="0.5" />
        </el-form-item>
        <el-form-item label="数量" prop="goodsAmount">
          <el-input-number v-model="formData.goodsAmount" :min="1" />
        </el-form-item>
        <el-form-item label="保价费(元)" prop="insuranceFee">
          <el-input-number v-model="formData.insuranceFee" :min="0" :precision="2" />
        </el-form-item>
        
        <el-divider>
          <span class="fee-info">预估运费: ¥{{ estimatedFee }}</span>
        </el-divider>
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">
            提交订单
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createOrderAPI } from '@/api'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const formData = reactive({
  senderName: '',
  senderPhone: '',
  senderAddress: '',
  receiverName: '',
  receiverPhone: '',
  receiverAddr: '',
  goodsType: '',
  goodsWeight: 1,
  goodsAmount: 1,
  insuranceFee: 0
})

const estimatedFee = computed(() => {
  return (formData.goodsWeight * 10 + formData.insuranceFee).toFixed(2)
})

const rules: FormRules = {
  senderName: [
    { required: true, message: '请输入寄件人姓名', trigger: 'blur' }
  ],
  senderPhone: [
    { required: true, message: '请输入寄件人电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  senderAddress: [
    { required: true, message: '请输入寄件人地址', trigger: 'blur' }
  ],
  receiverName: [
    { required: true, message: '请输入收件人姓名', trigger: 'blur' }
  ],
  receiverPhone: [
    { required: true, message: '请输入收件人电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  receiverAddr: [
    { required: true, message: '请输入收件人地址', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res: any = await createOrderAPI({
          sender_name: formData.senderName,
          sender_phone: formData.senderPhone,
          sender_address: formData.senderAddress,
          receiver_name: formData.receiverName,
          receiver_phone: formData.receiverPhone,
          receiver_addr: formData.receiverAddr,
          goods_type: formData.goodsType,
          goods_weight: formData.goodsWeight,
          goods_amount: formData.goodsAmount,
          insurance_fee: formData.insuranceFee
        })
        
        ElMessage.success(`订单创建成功！订单号：${res.data.order_no}`)
        router.push('/user/orders')
      } catch (error) {
        console.error('创建订单失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

const handleReset = () => {
  formRef.value?.resetFields()
}
</script>

<style scoped>
.create-order-container {
  padding: 0;
}

.fee-info {
  font-size: 16px;
  font-weight: bold;
  color: #f56c6c;
}
</style>
