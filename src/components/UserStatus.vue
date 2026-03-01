<template>
  <div class="flex items-center gap-3">
    <template v-if="isLoggedIn && user">
      <span class="text-sm text-foreground/70">你好, {{ user.nickname || user.username }}</span>
      <button 
        @click="handleLogout" 
        class="btn-ghost px-3 py-1.5 rounded-lg text-sm font-medium text-red-500 hover:text-red-600"
      >
        退出
      </button>
    </template>
    <template v-else>
      <a href="/login" class="btn-ghost px-3 py-1.5 rounded-lg text-sm font-medium">
        登录
      </a>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface User {
  id: number
  username: string
  nickname: string
  email: string
  role: string
}

const isLoggedIn = ref(false)
const user = ref<User | null>(null)

const API_BASE = import.meta.env.PUBLIC_API_BASE || 'http://localhost:8080'

// 检查登录状态
const checkAuth = async () => {
  const token = localStorage.getItem('token')
  if (!token) {
    isLoggedIn.value = false
    user.value = null
    return
  }

  try {
    const response = await fetch(`${API_BASE}/api/auth/me`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      user.value = data.user
      isLoggedIn.value = true
    } else {
      // Token 无效，清除
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      isLoggedIn.value = false
      user.value = null
    }
  } catch (error) {
    console.error('Failed to check auth:', error)
    isLoggedIn.value = false
    user.value = null
  }
}

// 退出登录
const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  isLoggedIn.value = false
  user.value = null
  // 刷新页面
  window.location.href = '/'
}

onMounted(() => {
  checkAuth()
})
</script>