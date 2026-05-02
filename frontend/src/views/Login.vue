<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h2>物流管理系统</h2>
        <p>Logistics Management System</p>
      </div>
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
      >
        <el-form-item prop="role">
          <el-radio-group v-model="loginForm.role">
            <el-radio value="user">用户</el-radio>
            <el-radio value="delivery">配送员</el-radio>
            <el-radio value="admin">管理员</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
        <div class="login-footer">
          <span v-if="loginForm.role === 'user'">
            还没有账号？<router-link to="/register">立即注册</router-link>
          </span>
          <span v-else>
            <el-text type="info">默认账号：admin / admin123</el-text>
          </span>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { loginAPI } from '@/api'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: '',
  role: 'user'
})

const loginRules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res: any = await loginAPI({
          username: loginForm.username,
          password: loginForm.password,
          role: loginForm.role
        })
        
        userStore.setToken(res.data.token)
        userStore.setRole(loginForm.role)
        userStore.setUserInfo(res.data.user)
        
        ElMessage.success('登录成功')
        
        const redirectMap: Record<string, string> = {
          user: '/user/dashboard',
          delivery: '/delivery/dashboard',
          admin: '/admin/dashboard'
        }
        router.push(redirectMap[loginForm.role] || '/user/dashboard')
      } catch (error) {
        console.error('登录失败:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h2 {
  margin: 0;
  color: #333;
  font-size: 24px;
}

.login-header p {
  margin: 10px 0 0;
  color: #999;
  font-size: 14px;
}

.login-form {
  margin-top: 20px;
}

.login-button {
  width: 100%;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
}

.login-footer a {
  color: #409eff;
  text-decoration: none;
}
</style>
