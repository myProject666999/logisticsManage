import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const userInfo = ref<any>(null)
  const role = ref<string | null>(localStorage.getItem('role'))

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const isUser = computed(() => role.value === 'user')
  const isDelivery = computed(() => role.value === 'delivery')
  const isAdmin = computed(() => role.value === 'admin')

  // 方法
  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setRole(newRole: string) {
    role.value = newRole
    localStorage.setItem('role', newRole)
  }

  function setUserInfo(info: any) {
    userInfo.value = info
    localStorage.setItem('userInfo', JSON.stringify(info))
  }

  function logout() {
    token.value = null
    userInfo.value = null
    role.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('role')
  }

  // 初始化
  function initFromStorage() {
    const storedUserInfo = localStorage.getItem('userInfo')
    if (storedUserInfo) {
      try {
        userInfo.value = JSON.parse(storedUserInfo)
      } catch {
        userInfo.value = null
      }
    }
  }

  return {
    token,
    userInfo,
    role,
    isLoggedIn,
    isUser,
    isDelivery,
    isAdmin,
    setToken,
    setRole,
    setUserInfo,
    logout,
    initFromStorage
  }
})
